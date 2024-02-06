package config

import "github.com/jaegertracing/jaeger/pkg/metrics"

var (
	metricsFactory metrics.Factory
)

func SetMetricsFactory(mf metrics.Factory) {
	metricsFactory = mf
}

func GetMetricsFactory() metrics.Factory {
	return metricsFactory
}
