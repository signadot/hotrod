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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/log"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/tracing"
)

// Client is a remote client that implements customer.Interface
type Client struct {
	tracer   opentracing.Tracer
	logger   log.Factory
	client   *tracing.HTTPClient
	hostPort string
}

// NewClient creates a new customer.Client
func NewClient(tracer opentracing.Tracer, logger log.Factory, hostPort string) *Client {
	return &Client{
		tracer: tracer,
		logger: logger,
		client: &tracing.HTTPClient{
			Client: &http.Client{Transport: &nethttp.Transport{}},
			Tracer: tracer,
		},
		hostPort: hostPort,
	}
}

// Get implements customer.Interface#Get as an RPC
func (c *Client) Get(ctx context.Context, customerID string) (*Customer, error) {
	c.logger.For(ctx).Info("Getting customer", zap.String("customer_id", customerID))

	url := fmt.Sprintf("http://"+c.hostPort+"/customer?customer=%s", customerID)
	fmt.Println(url)

	var customer Customer
	if err := c.client.GetJSON(ctx, "/customer", url, &customer); err != nil {
		return nil, err
	}
	return &customer, nil
}

func (c *Client) List(ctx context.Context) ([]Customer, error) {
	c.logger.For(ctx).Info("Getting customers")

	url := fmt.Sprintf("http://%s/customers", c.hostPort)
	fmt.Println(url)
	var customers []Customer
	if err := c.client.GetJSON(ctx, "/customers", url, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}
func (c *Client) Put(ctx context.Context, customer *Customer) error {
	c.logger.For(ctx).Info("PUT customer", zap.Int64("customer_id", customer.ID))

	url := fmt.Sprintf("http://%s/customer?customerID=%d", c.hostPort, customer.ID)
	fmt.Println(url)
	var outCustomer Customer
	if err := c.putJSON(ctx, "/customers", url, customer, outCustomer); err != nil {
		return err
	}
	return nil
}

func (c *Client) putJSON(ctx context.Context, ep, url string, reqIn, respOut interface{}) error {
	d, err := json.Marshal(reqIn)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(d))
	req.Header.Add("Content-type", "application/json")
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	req, ht := nethttp.TraceRequest(c.client.Tracer, req, nethttp.OperationName("HTTP PUT: "+ep))
	defer ht.Finish()

	res, err := c.client.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(respOut)
}
