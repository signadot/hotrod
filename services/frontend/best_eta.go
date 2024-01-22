// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package frontend

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/pool"
	"github.com/signadot/hotrod/services/config"
	"github.com/signadot/hotrod/services/driver"
	"github.com/signadot/hotrod/services/location"
	"github.com/signadot/hotrod/services/route"
)

type bestETA struct {
	location location.Interface
	driver   driver.Interface
	route    route.Interface
	pool     *pool.Pool
	logger   log.Factory
}

// Response contains ETA for a trip.
type Response struct {
	Driver string
	ETA    time.Duration
}

func newBestETA(tracer trace.TracerProvider, logger log.Factory, options ConfigOptions) *bestETA {
	return &bestETA{
		location: location.NewClient(
			tracer,
			logger.With(zap.String("component", "location_client")),
			options.LocationHostPort,
		),
		driver: driver.NewClient(
			tracer,
			logger.With(zap.String("component", "driver_client")),
			options.DriverHostPort,
		),
		route: route.NewClient(
			tracer,
			logger.With(zap.String("component", "route_client")),
			options.RouteHostPort,
		),
		pool:   pool.New(config.RouteWorkerPoolSize),
		logger: logger,
	}
}

func (eta *bestETA) Get(ctx context.Context, locationID int) (*Response, error) {
	location, err := eta.location.Get(ctx, locationID)
	if err != nil {
		return nil, err
	}
	eta.logger.For(ctx).Info("Found location", zap.Any("location", location))

	m, err := baggage.NewMember("location", location.Name)
	if err != nil {
		eta.logger.For(ctx).Error("cannot create baggage member", zap.Error(err))
	}
	bag := baggage.FromContext(ctx)
	bag, err = bag.SetMember(m)
	if err != nil {
		eta.logger.For(ctx).Error("cannot set baggage member", zap.Error(err))
	}
	ctx = baggage.ContextWithBaggage(ctx, bag)

	drivers, err := eta.driver.FindNearest(ctx, location.Coordinates)
	if err != nil {
		return nil, err
	}
	eta.logger.For(ctx).Info("Found drivers", zap.Any("drivers", drivers))

	results := eta.getRoutes(ctx, location, drivers)
	eta.logger.For(ctx).Info("Found routes", zap.Any("routes", results))

	resp := &Response{ETA: math.MaxInt64}
	for _, result := range results {
		if result.err != nil {
			return nil, err
		}
		if result.route.ETA < resp.ETA {
			resp.ETA = result.route.ETA
			resp.Driver = result.driver
		}
	}
	if resp.Driver == "" {
		return nil, errors.New("no routes found")
	}

	eta.logger.For(ctx).Info("Dispatch successful", zap.String("driver", resp.Driver), zap.String("eta", resp.ETA.String()))
	return resp, nil
}

type routeResult struct {
	driver string
	route  *route.Route
	err    error
}

// getRoutes calls Route service for each (location, driver) pair
func (eta *bestETA) getRoutes(ctx context.Context, pickUpLocation *location.Location, drivers []driver.Driver) []routeResult {
	results := make([]routeResult, 0, len(drivers))
	wg := sync.WaitGroup{}
	routesLock := sync.Mutex{}
	for _, dd := range drivers {
		wg.Add(1)
		driver := dd // capture loop var
		// Use worker pool to (potentially) execute requests in parallel
		eta.pool.Execute(func() {
			route, err := eta.route.FindRoute(ctx, driver.Coordinates, pickUpLocation.Coordinates)
			routesLock.Lock()
			results = append(results, routeResult{
				driver: driver.DriverID,
				route:  route,
				err:    err,
			})
			routesLock.Unlock()
			wg.Done()
		})
	}
	wg.Wait()
	return results
}
