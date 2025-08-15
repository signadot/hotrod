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
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/signadot/hotrod/internal/metrics/expvar"
	"github.com/signadot/hotrod/internal/metrics/prometheus"

	"github.com/jaegertracing/jaeger/pkg/metrics"

	"github.com/signadot/hotrod/pkg/config"
)

var (
	logger *zap.Logger
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hotrod",
	Short: "HotR.O.D. - A tracing demo application",
	Long:  `HotR.O.D. - A tracing demo application.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logger.Fatal("We bowled a googly", zap.Error(err))
		os.Exit(-1)
	}
}

func init() {
	addFlags(RootCmd)
	cobra.OnInitialize(onInitialize)
}

// onInitialize is called before the command is executed.
func onInitialize() {
	zapOptions := []zap.Option{
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	}
	if !verbose {
		zapOptions = append(zapOptions,
			zap.IncreaseLevel(zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l != zapcore.DebugLevel })),
		)
	}
	logger, _ = zap.NewDevelopment(zapOptions...)

	switch metricsBackend {
	case "expvar":
		config.SetMetricsFactory(expvar.NewFactory(10)) // 10 buckets for histograms
		logger.Info("*** Using expvar as metrics backend " + expvarDepr)
	case "prometheus":
		config.SetMetricsFactory(prometheus.New().Namespace(metrics.NSOptions{Name: "hotrod", Tags: nil}))
		logger.Info("Using Prometheus as metrics backend")
	default:
		logger.Fatal("unsupported metrics backend " + metricsBackend)
	}
}

func logError(logger *zap.Logger, err error) error {
	if err != nil {
		logger.Error("Error running command", zap.Error(err))
	}
	return err
}
