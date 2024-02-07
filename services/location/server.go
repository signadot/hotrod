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
	"fmt"
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
	hostPort       string
	tracerProvider trace.TracerProvider
	logger         log.Factory
	notification   notifications.Interface
	database       *database
}

// NewServer creates a new location.Server
func NewServer(hostPort string, logger log.Factory) *Server {
	// get a tracer provider for the location
	tracerProvider := tracing.InitOTEL("location", config.GetOtelExporterType(),
		config.GetMetricsFactory(), logger)

	return &Server{
		hostPort:       hostPort,
		tracerProvider: tracerProvider,
		logger:         logger,
		database:       newDatabase(logger),
		notification:   notifications.NewNotificationHandler(tracerProvider, logger),
	}
}

// Run starts the Location server
func (s *Server) Run() error {
	mux := s.createServeMux()
	s.logger.Bg().Info("Starting", zap.String("address", "http://"+s.hostPort))
	server := &http.Server{
		Addr:              s.hostPort,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux(false, s.tracerProvider, s.logger)
	mux.Handle("/locations", http.HandlerFunc(s.locations))
	mux.Handle("/location", http.HandlerFunc(s.location))
	mux.Handle("/healthz", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		s.logger.For(req.Context()).Info("/healthz")
		resp.Write([]byte("ok"))
	}))
	return mux
}

func (s *Server) locations(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) location(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	locationIDStr := r.Form.Get("locationID")
	if locationIDStr == "" {
		http.Error(w, "Missing required 'locationID' parameter", http.StatusBadRequest)
		return
	}
	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		http.Error(w, "Parameter 'locationID' is not an integer", http.StatusBadRequest)
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
