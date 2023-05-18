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
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/delay"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/log"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/tracing"
	"github.com/jaegertracing/jaeger/examples/hotrod/services/config"
)

// database simulates Customer repository implemented on top of an SQL database
type database struct {
	tracer    opentracing.Tracer
	logger    log.Factory
	customers map[string]*Customer
	lock      *tracing.Mutex
	db        *sqlx.DB
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
		Name:     "Bjarne Stroustrup",
		Location: "115,277",
	},
	{
		ID:       567,
		Name:     "Donald Knuth",
		Location: "211,653",
	},
	{
		ID:       392,
		Name:     "Tim Berners-Lee",
		Location: "577,322",
	},
	{
		ID:       731,
		Name:     "Vint Cerf",
		Location: "728,326",
	},
}

func newDatabase(tracer opentracing.Tracer, logger log.Factory) *database {
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
	dc.Addr = envDefault("MYSQL_HOST", "customer-db.hotrod.svc") +
		":" +
		envDefault("MYSQL_PORT", "3306")
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
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := d.tracer.StartSpan("SQL SELECT", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "mysql")
		// #nosec
		span.SetTag("sql.query", "SELECT id, name, location FROM customers")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
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

func (d *database) Get(ctx context.Context, customerID string) (*Customer, error) {
	d.logger.For(ctx).Info("Loading customer", zap.String("customer_id", customerID))

	// simulate opentracing instrumentation of an SQL query
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := d.tracer.StartSpan("SQL SELECT", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "mysql")
		// #nosec
		span.SetTag("sql.query", "SELECT * FROM customer WHERE customer_id="+customerID)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	cid, err := strconv.ParseInt(customerID, 10, 64)
	if err != nil {
		return nil, err
	}

	if !config.MySQLMutexDisabled {
		// simulate misconfigured connection pool that only gives one connection at a time
		d.lock.Lock(ctx)
		defer d.lock.Unlock()
	}

	// simulate RPC delay
	delay.Sleep(config.MySQLGetDelay, config.MySQLGetDelayStdDev)
	var c Customer
	if err := d.db.QueryRow("SELECT id, name, location from customers where ID= ?", cid).
		Scan(&c.ID, &c.Name, &c.Location); err != nil {
		return nil, err
	}
	return &c, nil
}
