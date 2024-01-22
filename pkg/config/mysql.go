package config

import "time"

func GetMySQLAddress() string {
	return EnvDefault("MYSQL_HOST", "location-db") +
		":" + EnvDefault("MYSQL_PORT", "3306")
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

var (
	// MySQLGetDelay is how long retrieving a location record takes.
	// Using large value mostly because I cannot click the button fast enough to cause a queue.
	MySQLGetDelay = 300 * time.Millisecond

	// MySQLGetDelayStdDev is standard deviation
	MySQLGetDelayStdDev = MySQLGetDelay / 10

	// MySQLMutexDisabled controls whether there is a mutex guarding db query execution.
	// When not disabled it simulates a misconfigured connection pool of size 1.
	MySQLMutexDisabled = false
)
