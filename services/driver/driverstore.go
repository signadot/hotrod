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

package driver

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/signadot/hotrod/pkg/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/signadot/hotrod/pkg/delay"
	"github.com/signadot/hotrod/pkg/log"
)

type driverStore struct {
	errorSimulator
	tracer trace.Tracer
	logger log.Factory
}

func newDriverStore(tracer trace.Tracer, logger log.Factory) *driverStore {
	return &driverStore{
		tracer: tracer,
		logger: logger,
		errorSimulator: errorSimulator{
			errorFreq: config.GetDriverStoreErrorFreq(),
		},
	}
}

// FindDriverIDs finds IDs of drivers who are near the location.
func (s *driverStore) FindDriverIDs(ctx context.Context) []string {
	ctx, span := s.tracer.Start(ctx, "FindDriverIDs", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	// simulate delay
	delay.Sleep(config.GetDriverStoreFindDelay(), config.GetDriverStoreFindStdDev())

	drivers := make([]string, 10)
	for i := range drivers {
		// #nosec
		drivers[i] = fmt.Sprintf("[[ %sN7%05dC%s ]]",
			config.GetDriverIDPrefix(), rand.Int()%100000, config.GetDriverIDSuffix())
	}
	s.logger.For(ctx).Info("Found drivers", zap.Strings("drivers", drivers))
	return drivers
}

func (s *driverStore) GetDriversLocation(ctx context.Context, driverIDs []string) ([]*Driver, error) {
	ctx, span := s.tracer.Start(ctx, "GetDriversLocation", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	drivers := make([]*Driver, len(driverIDs))
	for i, driverID := range driverIDs {
		var drv Driver
		var err error
		for i := 0; i < 3; i++ {
			drv, err = s.GetDriver(ctx, driverID)
			if err == nil {
				break
			}
			s.logger.For(ctx).Error("Retrying GetDriver after error",
				zap.Int("retry_no", i+1), zap.Error(err))
		}
		if err != nil {
			s.logger.For(ctx).Error("Failed to get driver after 3 attempts", zap.Error(err))
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		drivers[i] = &Driver{
			DriverID:    driverID,
			Coordinates: drv.Coordinates,
		}
	}
	return drivers, nil
}

// GetDriver returns driver and the current car location
func (s *driverStore) GetDriver(ctx context.Context, driverID string) (Driver, error) {
	ctx, span := s.tracer.Start(ctx, "GetDriver", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.
			Key("driver.id").
			String(driverID),
	)
	defer span.End()

	// simulate delay
	delay.Sleep(config.GetDriverStoreGetDelay(), config.GetDriverStoreGetStdDev())

	if err := s.checkError(); err != nil {
		// simulated error
		span.SetStatus(codes.Error, err.Error())
		s.logger.For(ctx).Error("error getting driver from ID", zap.String("driver_id", driverID), zap.Error(err))
		return Driver{}, err
	}

	return Driver{
		DriverID:    driverID,
		Coordinates: fmt.Sprintf("%d,%d", rand.Int()%1000, rand.Int()%1000),
	}, nil
}

var errTimeout = errors.New("driver store timeout")

type errorSimulator struct {
	sync.Mutex
	errorFreq      int
	countTillError int
}

func (es *errorSimulator) checkError() error {
	if es.errorFreq == 0 {
		return nil // errors disabled
	}
	es.Lock()
	es.countTillError--
	if es.countTillError > 0 {
		es.Unlock()
		return nil
	}
	es.countTillError = es.errorFreq
	es.Unlock()
	// add more delay for "timeout"
	delay.Sleep(config.GetDriverStoreGetDelay(), config.GetDriverStoreGetStdDev())
	return errTimeout
}
