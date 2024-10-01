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
)

// Location contains data about a location.
type Location struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%f,%f", c.Lat, c.Long)
}

// Interface exposed by the Location service.
type Interface interface {
	Get(ctx context.Context, locationID int) (*Location, error)
	List(ctx context.Context) ([]Location, error)
	Put(ctx context.Context, location *Location) error
}
