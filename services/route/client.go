package route

import (
	"context"
	"time"

	"github.com/signadot/hotrod/services/route/github.com/signadot/hotrod/services/route"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/signadot/hotrod/pkg/config"
	"github.com/signadot/hotrod/pkg/log"
)

// Client is a remote client that implements routes gRPC API
type Client struct {
	logger log.Factory
	client route.RoutesServiceClient
}

func NewClient(tracerProvider trace.TracerProvider, logger log.Factory) *Client {
	conn, err := grpc.Dial(
		config.GetRouteAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(otelgrpc.WithTracerProvider(tracerProvider))),
	)
	if err != nil {
		logger.Bg().Fatal("Cannot create gRPC connection", zap.Error(err))
	}

	client := route.NewRoutesServiceClient(conn)
	return &Client{
		logger: logger,
		client: client,
	}
}

func (c *Client) FindRoute(ctx context.Context, from, to string) (*Route, error) {
	c.logger.For(ctx).Info("Resolving route", zap.String("from", from), zap.String("to", from))
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	response, err := c.client.FindRoute(ctx, &route.FindRouteRequest{
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, err
	}
	return &Route{
		ETA: time.Duration(response.EtaSeconds) * time.Second,
	}, nil
}
