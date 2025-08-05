package config

import (
	"strings"
)

func GetNamespace() string {
	return EnvDefault("SIGNADOT_BASELINE_NAMESPACE", "hotrod")
}

func ExpandNamespace(v string) string {
	return strings.ReplaceAll(v, "<NAMESPACE>", GetNamespace())
}
