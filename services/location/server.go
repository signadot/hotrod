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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/signadot/hotrod/pkg/baggageutils"
	"github.com/signadot/hotrod/pkg/config"
	"github.com/signadot/hotrod/pkg/httperr"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/notifications"
	"github.com/signadot/hotrod/pkg/tracing"
)

// Server implements Location service
type Server struct {
	addr           string
	tracerProvider trace.TracerProvider
	logger         log.Factory
	notification   notifications.Interface
	database       *database
}

// NewServer creates new location.Server
func NewServer(logger log.Factory) *Server {
	// get tracer provider for the location
	tracerProvider := tracing.InitOTEL("location", logger)

	return &Server{
		addr:           net.JoinHostPort("0.0.0.0", config.GetLocationBindPort()),
		tracerProvider: tracerProvider,
		logger:         logger,
		database:       newDatabase(logger),
		notification:   notifications.NewNotificationHandler(tracerProvider, logger),
	}
}

// Run starts the Location server
func (s *Server) Run() error {
	mux := s.createServeMux()
	s.logger.Bg().Info("Starting", zap.String("address", "http://"+s.addr))
	server := &http.Server{
		Addr:              s.addr,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux(false, s.tracerProvider, s.logger)
	mux.Handle("GET /locations", http.HandlerFunc(s.listLocations))
	mux.Handle("GET /location", http.HandlerFunc(s.getLocation))
	mux.Handle("POST /location", http.HandlerFunc(s.createLocation))
	mux.Handle("DELETE /location", http.HandlerFunc(s.deleteLocation))
	mux.Handle("GET /healthz", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		s.logger.For(req.Context()).Info("/healthz")
		resp.Write([]byte("ok"))
	}))
	return mux
}

func (s *Server) listLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	locations, err := s.database.List(ctx)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot get locations", zap.Error(err))
		return
	}

	data, err := json.Marshal(locations)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot marshal response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) getLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	locationID, err := parseLocationID(r)
	if httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	response, err := s.database.Get(ctx, locationID)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("request failed", zap.Error(err))
		return
	}

	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot marshal response", zap.Error(err))
		return
	}

	// extract the request context
	reqContext, err := baggageutils.ExtractRequestContext(ctx)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot extract request context from baggage", zap.Error(err))
		return
	}
	if reqContext != nil {
		// send a notification
		s.notification.Store(ctx, &notifications.Notification{
			ID:        fmt.Sprintf("req-%d-location-resolve", reqContext.ID),
			Timestamp: time.Now(),
			Context:   s.notification.NotificationContext(reqContext, baggageutils.GetRoutingKey(ctx)),
			Body:      "Resolving locations",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) createLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))

	// decode the request body
	var loc Location
	err := json.NewDecoder(r.Body).Decode(&loc)
	if httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	// create the location
	loc.ID, err = s.database.Create(ctx, &loc)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("request failed", zap.Error(err))
		return
	}

	// marshal response
	data, err := json.Marshal(loc)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot marshal response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) deleteLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	locationID, err := parseLocationID(r)
	if httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	err = s.database.Delete(ctx, locationID)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("request failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func parseLocationID(r *http.Request) (int, error) {
	locationIDStr := r.Form.Get("locationID")
	if locationIDStr == "" {
		return 0, errors.New("missing required 'locationID' parameter")
	}
	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		return 0, errors.New("parameter 'locationID' is not an integer")
	}
	return locationID, nil
}
