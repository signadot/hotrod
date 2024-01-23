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

	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/services/frontend"
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
		var locationHost, routeHost string
		if baseDomain == "" {
			locationHost = "location"
			routeHost = "route"
		} else {
			if baseDomain == "localhost" || net.ParseIP(baseDomain) != nil {
				locationHost = baseDomain
				routeHost = baseDomain
			} else {
				locationHost = "location." + baseDomain
				routeHost = "route." + baseDomain
			}
		}
		if val := os.Getenv("FRONTEND_LOCATION_ADDR"); val != "" {
			options.LocationHostPort = val
		} else {
			options.LocationHostPort = net.JoinHostPort(locationHost, strconv.Itoa(locationPort))
		}
		if val := os.Getenv("FRONTEND_ROUTE_ADDR"); val != "" {
			options.RouteHostPort = val
		} else {
			options.RouteHostPort = net.JoinHostPort(routeHost, strconv.Itoa(routePort))
		}

		zapLogger := logger.With(zap.String("service", "frontend"))
		logger := log.NewFactory(zapLogger)
		server := frontend.NewServer(options, logger)
		return logError(zapLogger, server.Run())
	},
}

var options frontend.ConfigOptions

func init() {
	RootCmd.AddCommand(frontendCmd)
}
