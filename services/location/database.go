// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package location

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/signadot/hotrod/pkg/config"
	"github.com/signadot/hotrod/pkg/delay"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/tracing"
)

// database implements a Location repository on top of an SQL database
type database struct {
	tracer trace.Tracer
	logger log.Factory
	lock   *tracing.Mutex
	db     *sqlx.DB
}

const tableSchema = `
CREATE TABLE IF NOT EXISTS locations
(
    id bigint unsigned NOT NULL,
    name varchar(255) NOT NULL,
    coordinates varchar(255) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
`

var seed = []Location{
	{
		ID:          1,
		Name:        "My Home",
		Coordinates: "231,773",
	},
	{
		ID:          123,
		Name:        "Rachel's Floral Designs",
		Coordinates: "115,277",
	},
	{
		ID:          567,
		Name:        "Amazing Coffee Roasters",
		Coordinates: "211,653",
	},
	{
		ID:          392,
		Name:        "Trom Chocolatier",
		Coordinates: "577,322",
	},
	{
		ID:          731,
		Name:        "Japanese Desserts",
		Coordinates: "728,326",
	},
}

func newDatabase(logger log.Factory) *database {
	logger = logger.With(zap.String("component", "database"))

	db, err := sqlx.ConnectContext(context.TODO(), "mysql", driverConfig().FormatDSN())
	if err != nil {
		panic(err)
	}
	// Create the table if it doesn't already exist.
	_, err = db.Exec(tableSchema)
	if err != nil {
		panic(err)
	}
	count := 0
	if err := db.QueryRow("SELECT COUNT(*) from locations").Scan(&count); err != nil {
		panic(err)
	}
	if count == 0 {
		fmt.Println("seeding database")
		stmt, err := db.Prepare("INSERT into locations (id, name, coordinates) values (?, ?, ?)")
		if err != nil {
			panic(err)
		}
		for i := range seed {
			c := &seed[i]
			if _, err := stmt.Exec(c.ID, c.Name, c.Coordinates); err != nil {
				panic(err)
			}
		}
		stmt.Close()
	} else {
		fmt.Println("not seeding database")
	}

	return &database{
		tracer: tracing.InitOTEL("mysql", config.GetOtelExporterType(),
			config.GetMetricsFactory(), logger).Tracer("mysql"),
		logger: logger,
		lock: &tracing.Mutex{
			SessionBaggageKey: "request",
		},
		db: db,
	}
}

func driverConfig() *mysql.Config {
	dc := mysql.NewConfig()
	dc.Net = "tcp"
	dc.Addr = config.GetMySQLAddress()
	dc.DBName = config.GetMySQLDatabaseName()
	dc.User = config.GetMySQLUser()
	dc.Passwd = config.GetMySQLPassword()

	dc.Timeout = 60 * time.Second
	dc.InterpolateParams = true
	dc.ParseTime = true
	dc.Params = map[string]string{
		"time_zone": "'+00:00'",
	}
	return dc
}

func (d *database) List(ctx context.Context) ([]Location, error) {
	d.logger.For(ctx).Info("Loading locations", zap.String("location-id", "*"))
	// simulate opentracing instrumentation of an SQL query

	_, span := d.tracer.Start(ctx, "SQL SELECT", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		semconv.PeerServiceKey.String("mysql"),
		attribute.
			Key("sql.query").
			String("SELECT id, name, coordinates FROM locations"),
	)
	defer span.End()

	rows, err := d.db.Query("SELECT id, name, coordinates FROM locations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cs []Location
	for rows.Next() {
		c := Location{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Coordinates); err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
}

func (d *database) Put(ctx context.Context, location *Location) error {
	res, err := d.db.Exec("UPDATE locations set name = ?, coordinates = ? WHERE id = ?",
		location.Name, location.Coordinates, location.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("wrong number of rows on update: %d != 1", n)
	}
	return nil
}

func (d *database) Get(ctx context.Context, locationID int) (*Location, error) {
	d.logger.For(ctx).Info("Loading location", zap.Int("location_id", locationID))

	_, span := d.tracer.Start(ctx, "SQL SELECT", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		semconv.PeerServiceKey.String("mysql"),
		attribute.
			Key("sql.query").
			String(fmt.Sprintf("SELECT id, name, coordinates from locations where id = %d", locationID)),
	)
	defer span.End()

	// if !config.MySQLMutexDisabled {
	// 	// simulate misconfigured connection pool that only gives one connection at a time
	// 	d.lock.Lock(ctx)
	// 	defer d.lock.Unlock()
	// }

	// simulate RPC delay
	delay.Sleep(config.MySQLGetDelay, config.MySQLGetDelayStdDev)
	var c Location
	if err := d.db.QueryRow("SELECT id, name, coordinates from locations where id = ?", locationID).
		Scan(&c.ID, &c.Name, &c.Coordinates); err != nil {
		return nil, err
	}
	return &c, nil
}
