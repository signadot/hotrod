package config

import (
	"os"
	"time"
)

func GetRouteAddr() string {
	return EnvDefault("ROUTE_ADDR", "route:8083")
}

// how long a route calculation takes
func GetRouteCalcDelay() time.Duration {
	defaultDuration := 50 * time.Millisecond

	e := os.Getenv("ROUTE_CALC_DELAY")
	if e == "" {
		return defaultDuration
	}
	dur, err := time.ParseDuration(e)
	if err != nil {
		return defaultDuration
	}
	return dur
}

// GetRouteCalcStdDev calculates the standard deviation
func GetRouteCalcStdDev() time.Duration {
	return GetRouteCalcDelay() / 4
}
