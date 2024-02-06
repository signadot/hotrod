package config

func GetRedisAddr() string {
	return EnvDefault("REDIS_ADDR", "redis:6379")
}

func GetRedisPassword() string {
	return EnvDefault("REDIS_PASS", "")
}
