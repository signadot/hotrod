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

package cmd

import (
	"net"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/log"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/tracing"
	"github.com/jaegertracing/jaeger/examples/hotrod/services/frontend"
)

// frontendCmd represents the frontend command
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Starts Frontend service",
	Long:  `Starts Frontend service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		options.FrontendHostPort = net.JoinHostPort("0.0.0.0", strconv.Itoa(frontendPort))
		options.Basepath = basePath

		// Resolve services addresses
		var driverHost, customerHost, routeHost string
		if baseDomain == "" {
			driverHost = "driver"
			customerHost = "customer"
			routeHost = "route"
		} else {
			if baseDomain == "localhost" || net.ParseIP(baseDomain) != nil {
				driverHost = baseDomain
				customerHost = baseDomain
				routeHost = baseDomain
			} else {
				driverHost = "driver." + baseDomain
				customerHost = "customer." + baseDomain
				routeHost = "route." + baseDomain
			}
		}
		if val := os.Getenv("FRONTEND_DRIVER_ADDR"); val != "" {
			options.DriverHostPort = val
		} else {
			options.DriverHostPort = net.JoinHostPort(driverHost, strconv.Itoa(driverPort))
		}
		if val := os.Getenv("FRONTEND_CUSTOMER_ADDR"); val != "" {
			options.CustomerHostPort = val
		} else {
			options.CustomerHostPort = net.JoinHostPort(customerHost, strconv.Itoa(customerPort))
		}
		if val := os.Getenv("FRONTEND_ROUTE_ADDR"); val != "" {
			options.RouteHostPort = val
		} else {
			options.RouteHostPort = net.JoinHostPort(routeHost, strconv.Itoa(routePort))
		}

		zapLogger := logger.With(zap.String("service", "frontend"))
		logger := log.NewFactory(zapLogger)
		server := frontend.NewServer(
			options,
			tracing.Init("frontend", metricsFactory, logger),
			logger,
		)
		return logError(zapLogger, server.Run())
	},
}

var options frontend.ConfigOptions

func init() {
	RootCmd.AddCommand(frontendCmd)

}
