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
	"io/fs"
	"net/http"
	"path"
	"text/template"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/httperr"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/log"
	"github.com/jaegertracing/jaeger/examples/hotrod/pkg/tracing"
	"github.com/jaegertracing/jaeger/examples/hotrod/services/customer"
)

//go:embed web_assets/*
var webAssetsFS embed.FS

//go:embed templates/*
var tplFS embed.FS

// Server implements jaeger-demo-frontend service
type Server struct {
	hostPort   string
	tracer     opentracing.Tracer
	logger     log.Factory
	bestETA    *bestETA
	tplFS      fs.FS
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
	tplFS, err := fs.Sub(tplFS, "templates")
	_ = tplFS
	if err != nil {
		panic(err)
	}
	custClient := customer.NewClient(tracer, logger, options.CustomerHostPort)
	return &Server{
		hostPort:   options.FrontendHostPort,
		tracer:     tracer,
		logger:     logger,
		bestETA:    newBestETA(tracer, logger, options),
		tplFS:      tplFS,
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
	p := path.Join("/", s.basepath)
	ap := path.Join(p, "/web_assets")
	mux.Handle(ap+"/", http.StripPrefix(s.basepath, http.FileServer(http.FS(webAssetsFS))))
	dp := path.Join(p, "/dispatch")
	mux.Handle(dp, http.HandlerFunc(s.dispatch))
	mux.Handle(p, http.HandlerFunc(s.splash))
	return mux
}

func (s *Server) splash(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	t, err := template.ParseFS(s.tplFS, "*")
	if err != nil {
		httperr.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	cs, err := s.custClient.List(ctx)
	if err != nil {
		httperr.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	var rows [][]customer.Customer
	mod := 3
	for i, cust := range cs {
		if i%mod == 0 {
			rows = append(rows, []customer.Customer{})
		}
		ri := len(rows) - 1
		rows[ri] = append(rows[ri], cust)
	}
	if err := t.Execute(w, struct{ Rows [][]customer.Customer }{rows}); err != nil {
		httperr.HandleError(w, err, http.StatusInternalServerError)
		return
	}
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
