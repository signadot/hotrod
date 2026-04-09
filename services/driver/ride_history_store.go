package driver

import (
	"context"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type rideHistoryStore struct {
	db *sqlx.DB
}

const rideHistoryTableSchema = `
CREATE TABLE IF NOT EXISTS ride_history
(
	id bigint unsigned NOT NULL AUTO_INCREMENT,
	request_id bigint unsigned NOT NULL,
	pickup_location varchar(255) NOT NULL,
	dropoff_location varchar(255) NOT NULL,
	requested_at datetime(6) NOT NULL,
	driver_plate varchar(255) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY request_id (request_id),
	KEY requested_at (requested_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
`

func newRideHistoryStore() *rideHistoryStore {
	dsn := mysqlConfig().FormatDSN()
	connectTicker := time.NewTicker(time.Second / 3)
	defer connectTicker.Stop()

	var (
		db  *sqlx.DB
		err error
	)
	for {
		db, err = sqlx.ConnectContext(context.TODO(), "mysql", dsn)
		if err == nil {
			break
		}
		<-connectTicker.C
	}

	store := &rideHistoryStore{db: db}
	store.setupTable()
	return store
}

func mysqlConfig() *mysql.Config {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = envOrDefault("MYSQL_ADDR", "mysql:3306")
	cfg.DBName = envOrDefault("MYSQL_DBNAME", "location")
	cfg.User = envOrDefault("MYSQL_USER", "root")
	cfg.Passwd = os.Getenv("MYSQL_PASS")
	cfg.Timeout = 60 * time.Second
	cfg.InterpolateParams = true
	cfg.ParseTime = true
	cfg.Params = map[string]string{
		"time_zone": "'+00:00'",
	}
	return cfg
}

func envOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (s *rideHistoryStore) setupTable() {
	ticker := time.NewTicker(time.Second / 3)
	defer ticker.Stop()

	for {
		_, err := s.db.Exec(rideHistoryTableSchema)
		if err == nil {
			return
		}
		<-ticker.C
	}
}

func (s *rideHistoryStore) Store(ctx context.Context, requestID uint,
	pickupLocation, dropoffLocation string, requestedAt time.Time, driverPlate string) error {
	query := `
INSERT INTO ride_history (request_id, pickup_location, dropoff_location, requested_at, driver_plate)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	pickup_location = VALUES(pickup_location),
	dropoff_location = VALUES(dropoff_location),
	requested_at = VALUES(requested_at),
	driver_plate = VALUES(driver_plate)
`
	_, err := s.db.ExecContext(ctx, query, requestID, pickupLocation, dropoffLocation, requestedAt, driverPlate)
	return err
}
