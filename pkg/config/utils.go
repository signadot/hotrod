package config

import (
	"os"
)

func EnvDefault(key, value string) string {
	e := os.Getenv(key)
	if e == "" {
		return value
	}
	return e
}
