package config

func GetFrontendBindPort() string {
	return EnvDefault("FRONTEND_BIND_PORT", "8080")
}

func GetFrontendBasepath() string {
	return EnvDefault("FRONTEND_BASEPATH", "")
}
