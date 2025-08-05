package config

import (
	"os"
)

func GetOtelExporterType() string {
	return EnvDefault("OTEL_EXPORTER_TYPE", "otlp") // otlp or stdout
}

func InitOtelExporter() {
	addr := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if addr == "" {
		return
	}
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", ExpandNamespace(addr))
}
