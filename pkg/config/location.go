package config

func GetLocationBindPort() string {
	return EnvDefault("LOCATION_BIND_PORT", "8081")
}

func GetLocationAddr() string {
	return ExpandNamespace(EnvDefault("LOCATION_ADDR", "location:8081"))
}
