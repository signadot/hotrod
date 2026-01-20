package route

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/signadot/hotrod/pkg/baggageutils"
	"github.com/signadot/hotrod/pkg/config"
	"github.com/signadot/hotrod/pkg/delay"
	"github.com/signadot/hotrod/pkg/notifications"
	"github.com/signadot/hotrod/pkg/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/signadot/hotrod/pkg/log"
)

// Server implements jaeger-demo-frontend service
type Server struct {
	UnimplementedRoutesServiceServer
	addr           string
	tracerProvider trace.TracerProvider
	logger         log.Factory
	server         *grpc.Server
	notification   notifications.Interface
}

var _ RoutesServiceServer = (*Server)(nil)

// NewServer creates a new Server
func NewServer(logger log.Factory) *Server {
	// get a tracer provider for the route service
	tracerProvider := tracing.InitOTEL("route", logger)

	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(tracerProvider))),
	)

	return &Server{
		addr:           net.JoinHostPort("0.0.0.0", config.GetRouteBindPort()),
		tracerProvider: tracerProvider,
		logger:         logger,
		server:         server,
		notification:   notifications.NewNotificationHandler(tracerProvider, logger),
	}
}

// Run starts the Driver server
func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Bg().Fatal("Unable to create http listener", zap.Error(err))
	}
	RegisterRoutesServiceServer(s.server, s)
	reflection.Register(s.server)

	err = s.server.Serve(lis)
	if err != nil {
		s.logger.Bg().Fatal("Unable to start gRPC server", zap.Error(err))
	}
	return err
}

// FindRoute implements gRPC route interface
func (s *Server) FindRoute(ctx context.Context, req *FindRouteRequest) (*FindRouteResponse, error) {
	s.logger.For(ctx).Info("Finding route", zap.String("from", req.From), zap.String("to", req.To))

	// Simulate expensive calculation
	delay.Sleep(config.GetRouteCalcDelay(), config.GetRouteCalcStdDev())

	// Generate a random number between 3 and 45 with decimals
	eta := time.Duration((rand.Float64()*(45-3) + 3) * float64(time.Minute))
	if os.Getenv("FAST_ROUTE") != "" {
		eta = time.Second
	}
	// Round to the second
	eta = eta.Round(time.Second)

	// extract the request context
	reqContext, err := baggageutils.ExtractRequestContext(ctx)
	if err != nil {
		s.logger.For(ctx).Error("cannot extract request context from baggage", zap.Error(err))
		return nil, err
	}
	if reqContext != nil {
		// send a notification
		s.notification.Store(ctx, &notifications.Notification{
			ID:        fmt.Sprintf("req-%d-route-resolve", reqContext.ID),
			Timestamp: time.Now(),
			Context:   s.notification.NotificationContext(reqContext, baggageutils.GetRoutingKey(ctx)),
			Body:      "Resolving routes",
		})
	}

	// Generate random distance between 1 and 20 km
	distanceMeters := int32((rand.Float64()*(20-1) + 1) * 1000)

	return &FindRouteResponse{
		EtaSeconds:     int32(time.Duration(eta) / time.Second),
		DistanceMeters: distanceMeters,
	}, nil
}
