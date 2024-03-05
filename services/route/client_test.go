package route

import (
	"context"
	"os"
	"testing"

	"github.com/signadot/hotrod/pkg/config"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLocationClient(t *testing.T) {
	if os.Getenv("TEST_ROUTE_ADDR") == "" {
		t.Skip()
		return
	}
	ctx := context.Background()

	zapOptions := []zap.Option{
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
		zap.IncreaseLevel(zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l != zapcore.DebugLevel })),
	}
	logger, _ := zap.NewDevelopment(zapOptions...)
	zapLogger := logger.With(zap.String("service", "route"))
	l := log.NewFactory(zapLogger)

	// get a tracer provider for the frontend
	tracerProvider := tracing.InitOTEL("route", config.GetOtelExporterType(),
		config.GetMetricsFactory(), l)

	// get a location client
	routeClient := NewClient(tracerProvider, l, os.Getenv("TEST_ROUTE_ADDR"))
	route, err := routeClient.FindRoute(ctx, "a", "b")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("route: %+v", route)
}
