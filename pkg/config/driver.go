package config

import (
	"os"
	"strconv"
	"time"
)

func GetDriverIDPrefix() string {
	return os.Getenv("DRIVER_ID_PREFIX")
}

func GetDriverIDSuffix() string {
	return os.Getenv("DRIVER_ID_SUFFIX")
}

// how long finding closest drivers takes.
func GetDriverStoreFindDelay() time.Duration {
	defaultDuration := 20 * time.Millisecond

	e := os.Getenv("DRIVERSTORE_FIND_DELAY")
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
func GetDriverStoreFindStdDev() time.Duration {
	return GetDriverStoreFindDelay() / 4
}

// how long finding closest drivers takes.
func GetDriverStoreGetDelay() time.Duration {
	defaultDuration := 10 * time.Millisecond

	e := os.Getenv("DRIVERSTORE_GET_DELAY")
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
func GetDriverStoreGetStdDev() time.Duration {
	return GetDriverStoreGetDelay() / 4
}

// simulate timeout errors every x calls
func GetDriverStoreErrorFreq() int {
	defaultValue := 5

	e := os.Getenv("DRIVERSTORE_ERROR_FREQ")
	if e == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(e)
	if err != nil {
		return defaultValue
	}
	return val
}

// the size of the worker pool used to query `route` service.
func GetDriverWorkerPoolSize() int {
	defaultValue := 3

	e := os.Getenv("DRIVER_WORKER_POOL_SIZE")
	if e == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(e)
	if err != nil {
		return defaultValue
	}
	return val
}
