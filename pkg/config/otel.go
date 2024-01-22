package config

func GetOtelExporterType() string {
	return EnvDefault("OTEL_EXPORTER_TYPE", "otlp") // otlp or stdout
}
