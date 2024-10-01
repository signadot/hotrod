package driver

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/signadot/hotrod/services/location"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/signadot/hotrod/pkg/config"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/pool"
	"github.com/signadot/hotrod/services/route"
)

type bestETA struct {
	tracer trace.Tracer
	route  route.Interface
	pool   *pool.Pool
	logger log.Factory
}

// Response contains ETA for a trip.
type Response struct {
	DriverID string
	ETA      time.Duration
}

func newBestETA(tracerProvider trace.TracerProvider, tracer trace.Tracer, logger log.Factory) *bestETA {
	return &bestETA{
		tracer: tracer,
		route:  route.NewClient(tracerProvider, logger, config.GetRouteAddr()),
		pool:   pool.New(config.GetDriverWorkerPoolSize()),
		logger: logger,
	}
}

func (eta *bestETA) Get(ctx context.Context, dispatchReq *DispatchRequest,
	drivers []*Driver) (*Response, error) {
	ctx, span := eta.tracer.Start(ctx, "CalculateBestRoute", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	// get all routes from the drivers to the pick up location
	results := eta.getRoutes(ctx, dispatchReq.PickupLocation, drivers)
	eta.logger.For(ctx).Info("Found routes", zap.Any("routes", results))

	// search the one with the best ETA
	resp := &Response{ETA: math.MaxInt64}
	for _, result := range results {
		if result.err != nil {
			span.SetStatus(codes.Error, result.err.Error())
			return nil, result.err
		}
		if result.route.ETA < resp.ETA {
			resp.ETA = result.route.ETA
			resp.DriverID = result.driverID
		}
	}
	if resp.DriverID == "" {
		err := errors.New("no routes found")
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	eta.logger.For(ctx).Info("Dispatch successful",
		zap.String("driverID", resp.DriverID),
		zap.String("eta", resp.ETA.String()),
	)
	return resp, nil
}

type routeResult struct {
	driverID string
	route    *route.Route
	err      error
}

// getRoutes calls Route service for each pick location, driver location pair
func (eta *bestETA) getRoutes(ctx context.Context, pickupLoc *location.Location,
	drivers []*Driver) []routeResult {
	results := make([]routeResult, 0, len(drivers))
	wg := sync.WaitGroup{}
	routesLock := sync.Mutex{}

	for _, dd := range drivers {
		wg.Add(1)
		driver := dd // capture loop var
		// Use worker pool to (potentially) execute requests in parallel
		eta.pool.Execute(func() {
			route, err := eta.route.FindRoute(ctx, driver.Coordinates, pickupLoc.Coordinates.String())
			routesLock.Lock()
			results = append(results, routeResult{
				driverID: driver.DriverID,
				route:    route,
				err:      err,
			})
			routesLock.Unlock()
			wg.Done()
		})
	}
	wg.Wait()
	return results
}
