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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/tracing"
)

// Client is a remote client that implements location.Interface
type Client struct {
	logger   log.Factory
	client   *tracing.HTTPClient
	hostPort string
}

// NewClient creates a new location.Client
func NewClient(tracer trace.TracerProvider, logger log.Factory, hostPort string) *Client {
	return &Client{
		logger:   logger,
		client:   tracing.NewHTTPClient(tracer),
		hostPort: hostPort,
	}
}

// Get implements location.Interface#Get as an RPC
func (c *Client) Get(ctx context.Context, locationID int) (*Location, error) {
	c.logger.For(ctx).Info("Getting location", zap.Int("location_id", locationID))

	url := fmt.Sprintf("http://"+c.hostPort+"/location?locationID=%d", locationID)
	var location Location
	if err := c.client.GetJSON(ctx, "/location", url, &location); err != nil {
		return nil, err
	}
	return &location, nil
}

func (c *Client) List(ctx context.Context) ([]Location, error) {
	c.logger.For(ctx).Info("Getting locations")

	url := fmt.Sprintf("http://%s/locations", c.hostPort)
	fmt.Println(url)
	var locations []Location
	if err := c.client.GetJSON(ctx, "/locations", url, &locations); err != nil {
		return nil, err
	}
	return locations, nil
}

func (c *Client) Put(ctx context.Context, location *Location) error {
	c.logger.For(ctx).Info("PUT location", zap.Int64("location_id", location.ID))

	url := fmt.Sprintf("http://%s/location?locationID=%d", c.hostPort, location.ID)
	fmt.Println(url)
	var outLocation Location
	if err := c.putJSON(ctx, "/locations", url, location, outLocation); err != nil {
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
