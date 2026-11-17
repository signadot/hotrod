package driver

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/pool"
	"github.com/signadot/hotrod/services/location"
	"github.com/signadot/hotrod/services/route"
)

const skippedMsg = "route failed, driver skipped"

// fakeRoute answers FindRoute by the driver's coordinates (the `from` argument):
// an ETA, or an error.
type fakeRoute struct {
	etas map[string]time.Duration
	errs map[string]error
}

func (f *fakeRoute) FindRoute(ctx context.Context, from, to string) (*route.Route, error) {
	if err, ok := f.errs[from]; ok {
		return nil, err
	}
	if eta, ok := f.etas[from]; ok {
		return &route.Route{ETA: eta}, nil
	}
	return nil, fmt.Errorf("fakeRoute: no answer for %q", from)
}

type etaHarness struct {
	eta   *bestETA
	spans *tracetest.SpanRecorder
	logs  *observer.ObservedLogs
}

func newETAHarness(t *testing.T, fake *fakeRoute) *etaHarness {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })
	core, logs := observer.New(zapcore.DebugLevel)
	return &etaHarness{
		eta: &bestETA{
			tracer: tp.Tracer("test"),
			route:  fake,
			pool:   pool.New(3),
			logger: log.NewFactory(zap.New(core)),
		},
		spans: recorder,
		logs:  logs,
	}
}

func (h *etaHarness) span(t *testing.T) sdktrace.ReadOnlySpan {
	t.Helper()
	ended := h.spans.Ended()
	if len(ended) != 1 {
		t.Fatalf("got %d recorded spans, want 1", len(ended))
	}
	if ended[0].Name() != "CalculateBestRoute" {
		t.Fatalf("span name = %q, want CalculateBestRoute", ended[0].Name())
	}
	return ended[0]
}

// skipped returns the driver_id of each "route failed, driver skipped" Error
// log entry, sorted, and checks each carries an error field.
func (h *etaHarness) skipped(t *testing.T) []string {
	t.Helper()
	var ids []string
	for _, e := range h.logs.FilterMessage(skippedMsg).All() {
		if e.Level != zapcore.ErrorLevel {
			t.Errorf("%q logged at %v, want error", skippedMsg, e.Level)
		}
		ctx := e.ContextMap()
		if _, ok := ctx["error"]; !ok {
			t.Errorf("%q entry has no error field: %v", skippedMsg, ctx)
		}
		id, _ := ctx["driver_id"].(string)
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func assertCounts(t *testing.T, span sdktrace.ReadOnlySpan, asked, answered, skipped int) {
	t.Helper()
	want := map[attribute.Key]int64{
		"drivers.asked":    int64(asked),
		"drivers.answered": int64(answered),
		"drivers.skipped":  int64(skipped),
	}
	got := map[attribute.Key]int64{}
	for _, kv := range span.Attributes() {
		if _, ok := want[kv.Key]; ok {
			got[kv.Key] = kv.Value.AsInt64()
		}
	}
	for k, v := range want {
		g, ok := got[k]
		if !ok {
			t.Errorf("span attribute %s missing", k)
		} else if g != v {
			t.Errorf("span attribute %s = %d, want %d", k, g, v)
		}
	}
}

func countEvents(span sdktrace.ReadOnlySpan, name string) int {
	n := 0
	for _, e := range span.Events() {
		if e.Name == name {
			n++
		}
	}
	return n
}

func drivers(ids ...string) []*Driver {
	ds := make([]*Driver, 0, len(ids))
	for _, id := range ids {
		ds = append(ds, &Driver{DriverID: id, Coordinates: "loc-" + id})
	}
	return ds
}

func dispatchRequest() *DispatchRequest {
	return &DispatchRequest{PickupLocation: &location.Location{Coordinates: "pickup"}}
}

func TestBestETASkipsFailedRoutes(t *testing.T) {
	h := newETAHarness(t, &fakeRoute{
		etas: map[string]time.Duration{"loc-d2": 10 * time.Minute, "loc-d4": 5 * time.Minute},
		errs: map[string]error{"loc-d1": errors.New("route d1 down"), "loc-d3": errors.New("route d3 down")},
	})

	resp, err := h.eta.Get(context.Background(), dispatchRequest(), drivers("d1", "d2", "d3", "d4"))
	if err != nil {
		t.Fatalf("Get returned error %v, want the fastest answering driver", err)
	}
	if resp == nil || *resp != (Response{DriverID: "d4", ETA: 5 * time.Minute}) {
		t.Fatalf("Get = %+v, want {d4 5m}", resp)
	}
	if got := h.skipped(t); strings.Join(got, ",") != "d1,d3" {
		t.Errorf("skipped drivers logged = %v, want [d1 d3]", got)
	}
	span := h.span(t)
	if span.Status().Code != codes.Unset {
		t.Errorf("span status = %v, want Unset", span.Status().Code)
	}
	assertCounts(t, span, 4, 2, 2)
	if n := countEvents(span, skippedMsg); n != 2 {
		t.Errorf("span has %d %q events, want 2", n, skippedMsg)
	}
}

func TestBestETAOneAnswerIsEnough(t *testing.T) {
	h := newETAHarness(t, &fakeRoute{
		etas: map[string]time.Duration{"loc-d3": 7 * time.Minute},
		errs: map[string]error{"loc-d1": errors.New("down"), "loc-d2": errors.New("down")},
	})

	resp, err := h.eta.Get(context.Background(), dispatchRequest(), drivers("d1", "d2", "d3"))
	if err != nil {
		t.Fatalf("Get returned error %v, want d3", err)
	}
	if resp == nil || *resp != (Response{DriverID: "d3", ETA: 7 * time.Minute}) {
		t.Fatalf("Get = %+v, want {d3 7m}", resp)
	}
	assertCounts(t, h.span(t), 3, 1, 2)
}

func TestBestETAAllRoutesFail(t *testing.T) {
	sentinels := []error{errors.New("e1"), errors.New("e2"), errors.New("e3")}
	h := newETAHarness(t, &fakeRoute{
		errs: map[string]error{"loc-d1": sentinels[0], "loc-d2": sentinels[1], "loc-d3": sentinels[2]},
	})

	resp, err := h.eta.Get(context.Background(), dispatchRequest(), drivers("d1", "d2", "d3"))
	if resp != nil {
		t.Errorf("Get response = %+v, want nil", resp)
	}
	if err == nil {
		t.Fatal("Get returned no error, want no routes found")
	}
	if !strings.HasPrefix(err.Error(), "no routes found") || !strings.Contains(err.Error(), "0 of 3") {
		t.Errorf("error = %q, want it to start with \"no routes found\" and say \"0 of 3\"", err)
	}
	wrapped := false
	for _, s := range sentinels {
		if errors.Is(err, s) {
			wrapped = true
		}
	}
	if !wrapped {
		t.Errorf("error %q wraps none of the route errors", err)
	}
	span := h.span(t)
	if span.Status().Code != codes.Error {
		t.Errorf("span status = %v, want Error", span.Status().Code)
	}
	assertCounts(t, span, 3, 0, 3)
	if got := h.skipped(t); len(got) != 3 {
		t.Errorf("got %d %q log entries, want 3", len(got), skippedMsg)
	}
}

func TestBestETAAllAnswer(t *testing.T) {
	h := newETAHarness(t, &fakeRoute{
		etas: map[string]time.Duration{"loc-d1": 9 * time.Minute, "loc-d2": 2 * time.Minute, "loc-d3": 4 * time.Minute},
	})

	resp, err := h.eta.Get(context.Background(), dispatchRequest(), drivers("d1", "d2", "d3"))
	if err != nil {
		t.Fatalf("Get returned error %v", err)
	}
	if resp == nil || *resp != (Response{DriverID: "d2", ETA: 2 * time.Minute}) {
		t.Fatalf("Get = %+v, want {d2 2m}", resp)
	}
	if got := h.skipped(t); len(got) != 0 {
		t.Errorf("skipped drivers logged = %v, want none", got)
	}
	span := h.span(t)
	if span.Status().Code != codes.Unset {
		t.Errorf("span status = %v, want Unset", span.Status().Code)
	}
	assertCounts(t, span, 3, 3, 0)
}

func TestBestETANoDrivers(t *testing.T) {
	h := newETAHarness(t, &fakeRoute{})

	resp, err := h.eta.Get(context.Background(), dispatchRequest(), nil)
	if resp != nil {
		t.Errorf("Get response = %+v, want nil", resp)
	}
	if err == nil {
		t.Fatal("Get returned no error, want no routes found")
	}
	if !strings.HasPrefix(err.Error(), "no routes found") || !strings.Contains(err.Error(), "0 of 0") {
		t.Errorf("error = %q, want it to start with \"no routes found\" and say \"0 of 0\"", err)
	}
	if errors.Unwrap(err) != nil {
		t.Errorf("error wraps %v, want nothing", errors.Unwrap(err))
	}
	if h.span(t).Status().Code != codes.Error {
		t.Errorf("span status = %v, want Error", h.span(t).Status().Code)
	}
}
