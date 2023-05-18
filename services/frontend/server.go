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

package frontend

import (
	"embed"
	"encoding/json"
	"net/http"
	"path"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/httperr"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/log"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/tracing"
	"github.com/jaegertracing/jaeger/examples/hotrod/services/customer"
)

//go:embed client/build/*
var reactAppFS embed.FS

// Server implements jaeger-demo-frontend service
type Server struct {
	hostPort   string
	tracer     opentracing.Tracer
	logger     log.Factory
	bestETA    *bestETA
	basepath   string
	custClient *customer.Client
}

// ConfigOptions used to make sure service clients
// can find correct server ports
type ConfigOptions struct {
	FrontendHostPort string
	DriverHostPort   string
	CustomerHostPort string
	RouteHostPort    string
	Basepath         string
}

// NewServer creates a new frontend.Server
func NewServer(options ConfigOptions, tracer opentracing.Tracer, logger log.Factory) *Server {
	custClient := customer.NewClient(tracer, logger, options.CustomerHostPort)
	return &Server{
		hostPort:   options.FrontendHostPort,
		tracer:     tracer,
		logger:     logger,
		bestETA:    newBestETA(tracer, logger, options),
		basepath:   options.Basepath,
		custClient: custClient,
	}
}

// Run starts the frontend server
func (s *Server) Run() error {
	mux := s.createServeMux()
	s.logger.Bg().Info("Starting", zap.String("address", "http://"+path.Join(s.hostPort, s.basepath)))
	return http.ListenAndServe(s.hostPort, mux)
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux(s.tracer)
	p := path.Join("/api", s.basepath)

	// handle the built React application files
	reactAppPath := path.Join(p, "/")
	mux.Handle(reactAppPath, http.StripPrefix(s.basepath, http.FileServer(http.FS(reactAppFS))))

	// handle your API endpoints
	cp := path.Join(p, "/customers")
	mux.Handle(cp, http.HandlerFunc(s.customers))

	// handle your API endpoints
	dp := path.Join(p, "/dispatch")
	mux.Handle(dp, http.HandlerFunc(s.dispatch))

	// handle data for the splash page
	sp := path.Join(p, "/splash")
	mux.Handle(sp, http.HandlerFunc(s.splash))

	return mux
}

func (s *Server) customers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	cs, err := s.custClient.List(ctx)
	if err != nil {
		httperr.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	// Respond with JSON data
	jsonResponse, err := json.Marshal(cs)
	if err != nil {
		httperr.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (s *Server) splash(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	cs, err := s.custClient.List(ctx)
	if err != nil {
		httperr.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	var rows [][]customer.Customer
	mod := 4
	for i, cust := range cs {
		if i%mod == 0 {
			rows = append(rows, []customer.Customer{})
		}
		ri := len(rows) - 1
		rows[ri] = append(rows[ri], cust)
	}
	// Respond with JSON data
	jsonResponse, err := json.Marshal(rows)
	if err != nil {
		httperr.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (s *Server) dispatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	customerID := r.Form.Get("customer")
	if customerID == "" {
		http.Error(w, "Missing required 'customer' parameter", http.StatusBadRequest)
		return
	}
	// TODO distinguish between user errors (such as invalid customer ID) and server failures
	response, err := s.bestETA.Get(ctx, customerID)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("request failed", zap.Error(err))
		return
	}

	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot marshal response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
