package config

import (
	"os"
	"time"
)

func GetMySQLAddr() string {
	return EnvDefault("MYSQL_ADDR", "mysql:3306")
}

func GetMySQLUser() string {
	return EnvDefault("MYSQL_USER", "root")
}

func GetMySQLPassword() string {
	return EnvDefault("MYSQL_PASS", "")
}

func GetMySQLDatabaseName() string {
	return EnvDefault("MYSQL_DBNAME", "location")
}

// how long retrieving a location record takes.
func GetMySQLGetDelay() time.Duration {
	defaultDuration := 300 * time.Millisecond

	e := os.Getenv("MYSQL_GET_DELAY")
	if e == "" {
		return defaultDuration
	}
	dur, err := time.ParseDuration(e)
	if err != nil {
		return defaultDuration
	}
	return dur
}

// the standard deviation
func GetMySQLGetDelayStdDev() time.Duration {
	return GetMySQLGetDelay() / 10
}
