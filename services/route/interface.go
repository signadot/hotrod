package route

import (
	"context"
	"time"
)

// Route describes a route.
type Route struct {
	ETA time.Duration
}

// Interface exposed by the Route service.
type Interface interface {
	FindRoute(ctx context.Context, from, to string) (*Route, error)
}
