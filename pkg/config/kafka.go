package config

func GetKafkaBrokerAddr() string {
	return ExpandNamespace(EnvDefault("KAFKA_BROKER_ADDR", "kafka-headless:9092"))
}

func GetKafkaBrokers() []string {
	return []string{GetKafkaBrokerAddr()}
}
