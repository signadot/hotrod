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

package route

import (
	"context"
	"expvar"
	"time"

	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/tracing"
)

var (
	routeCalcByCustomer = expvar.NewMap("route.calc.by.customer.sec")
	routeCalcBySession  = expvar.NewMap("route.calc.by.session.sec")
)

var stats = []struct {
	expvar     *expvar.Map
	baggageKey string
}{
	{
		expvar:     routeCalcByCustomer,
		baggageKey: "customer",
	},
	{
		expvar:     routeCalcBySession,
		baggageKey: "session",
	},
}

func updateCalcStats(ctx context.Context, delay time.Duration) {
	delaySec := float64(delay/time.Millisecond) / 1000.0
	for _, s := range stats {
		key := tracing.BaggageItem(ctx, s.baggageKey)
		if key != "" {
			s.expvar.AddFloat(key, delaySec)
		}
	}
}
