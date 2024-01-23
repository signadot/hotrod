// Copyright (c) 2023 The Jaeger Authors.
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
	"github.com/spf13/cobra"
)

var (
	metricsBackend string
	verbose        bool

	locationPort int
	frontendPort int
	routePort    int

	basePath   string
	baseDomain string
)

const expvarDepr = "(deprecated, will be removed after 2024-01-01 or in release v1.53.0, whichever is later) "

// used by root command
func addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&metricsBackend, "metrics", "m", "prometheus", expvarDepr+"Metrics backend (expvar|prometheus). ")

	// Add flags to choose ports for services
	cmd.PersistentFlags().IntVarP(&locationPort, "location-service-port", "c", 8081, "Port for location service")
	cmd.PersistentFlags().IntVarP(&frontendPort, "frontend-service-port", "f", 8080, "Port for frontend service")
	cmd.PersistentFlags().IntVarP(&routePort, "route-service-port", "r", 8083, "Port for routing service")

	// Flag for serving frontend at custom basepath url
	cmd.PersistentFlags().StringVarP(&basePath, "basepath", "b", "", `Basepath for frontend service(default "/")`)
	cmd.PersistentFlags().StringVarP(&baseDomain, "basedomain", "B", "", `Base domain for accessing hotrod services`)

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enables debug logging")
}
