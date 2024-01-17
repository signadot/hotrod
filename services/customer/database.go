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

package customer

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/signadot/hotrod/pkg/delay"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/tracing"
	"github.com/signadot/hotrod/services/config"
)

// database simulates Customer repository implemented on top of an SQL database
type database struct {
	tracer trace.Tracer
	logger log.Factory
	lock   *tracing.Mutex
	db     *sqlx.DB
}

const tableSchema = `
CREATE TABLE IF NOT EXISTS customers
(
    id bigint unsigned NOT NULL,

    name varchar(255) NOT NULL,

    location varchar(255) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
`

var seed = []Customer{
	{
		ID:       123,
		Name:     "Rachel's Floral Designs",
		Location: "115,277",
	},
	{
		ID:       567,
		Name:     "Amazing Coffee Roasters",
		Location: "211,653",
	},
	{
		ID:       392,
		Name:     "Trom Chocolatier",
		Location: "577,322",
	},
	{
		ID:       731,
		Name:     "Japanese Desserts",
		Location: "728,326",
	},
}

func newDatabase(tracer trace.Tracer, logger log.Factory) *database {
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
	if err := db.QueryRow("SELECT COUNT(*) from customers").Scan(&count); err != nil {
		panic(err)
	}
	if count == 0 {
		fmt.Println("seeding database")
		stmt, err := db.Prepare("INSERT into customers (id, name, location) values (?, ?, ?)")
		if err != nil {
			panic(err)
		}
		for i := range seed {
			c := &seed[i]
			if _, err := stmt.Exec(c.ID, c.Name, c.Location); err != nil {
				panic(err)
			}
		}
		stmt.Close()
	} else {
		fmt.Println("not seeding database")
	}

	return &database{
		tracer: tracer,
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
	dc.Addr = envDefault("MYSQL_HOST", "customer-db") +
		":" + envDefault("MYSQL_PORT", "3306")
	dc.DBName = "customer"
	dc.User = "root"
	dc.Passwd = envDefault("MYSQL_ROOT_PASSWORD", "abc")

	dc.Timeout = 60 * time.Second
	dc.InterpolateParams = true
	dc.ParseTime = true
	dc.Params = map[string]string{
		"time_zone": "'+00:00'",
	}
	return dc
}

func envDefault(key, value string) string {
	e := os.Getenv(key)
	if e == "" {
		return value
	}
	return e
}

func (d *database) List(ctx context.Context) ([]Customer, error) {
	d.logger.For(ctx).Info("Loading customers", zap.String("customer-id", "*"))
	// simulate opentracing instrumentation of an SQL query

	ctx, span := d.tracer.Start(ctx, "SQL SELECT", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		semconv.PeerServiceKey.String("mysql"),
		attribute.
			Key("sql.query").
			String("SELECT id, name, location FROM customers"),
	)
	defer span.End()

	rows, err := d.db.Query("SELECT id, name, location FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cs []Customer
	for rows.Next() {
		c := Customer{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Location); err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
}

func (d *database) Put(ctx context.Context, customer *Customer) error {
	res, err := d.db.Exec("UPDATE customers set name = ?, location = ? WHERE id = ?",
		customer.Name, customer.Location, customer.ID)
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

func (d *database) Get(ctx context.Context, customerID int) (*Customer, error) {
	d.logger.For(ctx).Info("Loading customer", zap.Int("customer_id", customerID))

	ctx, span := d.tracer.Start(ctx, "SQL SELECT", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		semconv.PeerServiceKey.String("mysql"),
		attribute.
			Key("sql.query").
			String(fmt.Sprintf("SELECT * FROM customer WHERE customer_id=%d", customerID)),
	)
	defer span.End()

	// if !config.MySQLMutexDisabled {
	// 	// simulate misconfigured connection pool that only gives one connection at a time
	// 	d.lock.Lock(ctx)
	// 	defer d.lock.Unlock()
	// }

	// simulate RPC delay
	delay.Sleep(config.MySQLGetDelay, config.MySQLGetDelayStdDev)
	var c Customer
	if err := d.db.QueryRow("SELECT id, name, location from customers where ID= ?", customerID).
		Scan(&c.ID, &c.Name, &c.Location); err != nil {
		return nil, err
	}
	return &c, nil
}
