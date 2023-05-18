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
	"strconv"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/delay"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/log"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/tracing"
	"github.com/jaegertracing/jaeger/examples/hotrod/services/config"
)

var names = []string{"Silver Head", "Sparky", "Andy Roid", "Combot", "Artron", "Brobot", "Servo", "Cyborgan", "Robyte", "Lyftr", "Machina", "Dee Royd"}

const imageBucketPath = "https://signadot-public.s3.us-west-2.amazonaws.com/bots/"

// Redis is a simulator of remote Redis cache
type Redis struct {
	tracer opentracing.Tracer // simulate redis as a separate process
	logger log.Factory
	errorSimulator
}

func newRedis(metricsFactory metrics.Factory, logger log.Factory) *Redis {
	return &Redis{
		tracer: tracing.Init("redis", metricsFactory, logger),
		logger: logger,
	}
}

// FindDriverIDs finds IDs of drivers who are near the location.
func (r *Redis) FindDriverIDs(ctx context.Context, location string) []string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := r.tracer.StartSpan("FindDriverIDs", opentracing.ChildOf(span.Context()))
		span.SetTag("param.location", location)
		ext.SpanKindRPCClient.Set(span)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	// simulate RPC delay
	delay.Sleep(config.RedisFindDelay, config.RedisFindDelayStdDev)

	drivers := make([]string, len(names))
	for i := range drivers {
		// #nosec
		drivers[i] = strconv.Itoa(i + 1)
	}
	r.logger.For(ctx).Info("Found drivers", zap.Strings("drivers", drivers))
	return drivers
}

// GetDriver returns driver and the current car location
func (r *Redis) GetDriver(ctx context.Context, driverID string) (Driver, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := r.tracer.StartSpan("GetDriver", opentracing.ChildOf(span.Context()))
		span.SetTag("param.driverID", driverID)
		ext.SpanKindRPCClient.Set(span)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	// simulate RPC delay
	delay.Sleep(config.RedisGetDelay, config.RedisGetDelayStdDev)
	if err := r.checkError(); err != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			ext.Error.Set(span, true)
		}
		r.logger.For(ctx).Error("redis timeout", zap.String("driver_id", driverID), zap.Error(err))
		return Driver{}, err
	}

	// #nosec
	idx, err := strconv.Atoi(driverID)
	if err != nil {
		panic("driverID expected to be a number")
	}
	return Driver{
		DriverID: driverID,
		Name:     names[idx-1],
		Location: fmt.Sprintf("%d,%d", rand.Int()%1000, rand.Int()%1000),
		ImageURL: fmt.Sprintf("%sbot%s.png", imageBucketPath, driverID),
	}, nil
}

var errTimeout = errors.New("redis timeout")

type errorSimulator struct {
	sync.Mutex
	countTillError int
}

func (es *errorSimulator) checkError() error {
	es.Lock()
	es.countTillError--
	if es.countTillError > 0 {
		es.Unlock()
		return nil
	}
	es.countTillError = 5
	es.Unlock()
	delay.Sleep(2*config.RedisGetDelay, 0) // add more delay for "timeout"
	return errTimeout
}
