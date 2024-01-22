package config

func GetKafkaBrokers() []string {
	broker := EnvDefault("KAFKA_BROKER", "kafka-headless:9092")
	return []string{broker}
}
