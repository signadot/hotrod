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
	"expvar"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/signadot/hotrod/pkg/baggageutils"
	"github.com/signadot/hotrod/pkg/config"
	"github.com/signadot/hotrod/pkg/httperr"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/notifications"
	"github.com/signadot/hotrod/pkg/tracing"
	"github.com/signadot/hotrod/services/location"
)

//go:embed web_assets/*
var assetFS embed.FS

//go:embed templates/*
var tplFS embed.FS

// Server implements hotrod-frontend service
type Server struct {
	hostPort string
	basepath string
	jaegerUI string

	tracer       trace.TracerProvider
	logger       log.Factory
	tplFS        fs.FS
	location     location.Interface
	notification notifications.Interface
	dispatcher   *dispatcher
}

// ConfigOptions used to make sure service clients
// can find correct server ports
type ConfigOptions struct {
	FrontendHostPort string
	LocationHostPort string
	RouteHostPort    string
	Basepath         string
}

// NewServer creates a new frontend.Server
func NewServer(options ConfigOptions, logger log.Factory) *Server {
	// load templates
	tplFS, err := fs.Sub(tplFS, "templates")
	_ = tplFS
	if err != nil {
		panic(err)
	}

	// get a tracer provider for the frontend
	tracerProvider := tracing.InitOTEL("frontend", config.GetOtelExporterType(),
		config.GetMetricsFactory(), logger)

	// get a location client
	locationClient := location.NewClient(tracerProvider, logger, options.LocationHostPort)

	// get a notification handler
	notificationHandler := notifications.NewNotificationHandler(tracerProvider, logger)

	// get a dispatcher
	dispatcher := newDispatcher(tracerProvider, logger, locationClient)

	return &Server{
		hostPort: options.FrontendHostPort,
		basepath: options.Basepath,

		tracer:       tracerProvider,
		logger:       logger,
		tplFS:        tplFS,
		location:     locationClient,
		notification: notificationHandler,
		dispatcher:   dispatcher,
	}
}

// Run starts the frontend server
func (s *Server) Run() error {
	mux := s.createServeMux()
	s.logger.Bg().Info("Starting", zap.String("address", "http://"+path.Join(s.hostPort, s.basepath)))
	server := &http.Server{
		Addr:              s.hostPort,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux(true, s.tracer, s.logger)
	p := path.Join("/", s.basepath)
	mux.Handle(path.Join(p, "/dispatch"), http.HandlerFunc(s.dispatch))
	mux.Handle(path.Join(p, "/notifications"), http.HandlerFunc(s.notifications))
	mux.Handle(path.Join(p, "/debug/vars"), expvar.Handler())       // expvar
	mux.Handle(path.Join(p, "/metrics"), promhttp.Handler())        // Prometheus
	mux.Handle(path.Join(p, "/splash"), http.HandlerFunc(s.splash)) // Prometheus
	staticFileServer := http.FileServer(http.FS(assetFS))

	mux.Handle("/web_assets/", http.StripPrefix("", staticFileServer))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			indexContent, err := assetFS.ReadFile("web_assets/index.html")
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(indexContent)
		} else {
			staticFileServer.ServeHTTP(w, r)
		}
	}))
	mux.Handle(path.Join(p, "/healthz"), http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		s.logger.For(req.Context()).Info("/healthz")
		resp.Write([]byte("ok"))
	}))
	return mux
}

func (s *Server) splash(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))

	// get all stored locations
	locations, err := s.location.List(ctx)
	if err != nil {
		httperr.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// load data for the template
	data := struct {
		Locations   []location.Location
		TitleSuffix string
	}{locations, os.Getenv("FRONTEND_TITLE_SUFFIX")}

	s.writeResponse(data, w, r)
}

func (s *Server) dispatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))

	// decode request body
	var dispatchReq DispatchRequest
	err := json.NewDecoder(r.Body).Decode(&dispatchReq)
	if httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("error parsing request body", zap.Error(err))
		return
	}

	// get request context
	reqContext := &notifications.RequestContext{
		ID:                dispatchReq.RequestID,
		SessionID:         dispatchReq.SessionID,
		PickupLocationID:  dispatchReq.PickupLocationID,
		DropoffLocationID: dispatchReq.DropoffLocationID,
	}

	// inject the request context into baggage and update the current context
	bag, err := baggageutils.InjectRequestContext(ctx, reqContext)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot inject request context into baggage", zap.Error(err))
		return
	}
	ctx = baggage.ContextWithBaggage(ctx, *bag)

	// send frontend-dispatching-driver notification
	s.logger.For(ctx).Info("Dispatching driver", zap.Any("request", dispatchReq))
	s.notification.Store(ctx, &notifications.Notification{
		ID:        fmt.Sprintf("req-%d-frontend-dispatching-driver", dispatchReq.RequestID),
		Timestamp: time.Now(),
		Context:   s.notification.NotificationContext(reqContext, baggageutils.GetRoutingKey(ctx)),
		Body:      "Processing dispatch driver request",
	})

	// resolve locations
	pickupLoc, dropoffLoc, err := s.dispatcher.ResolveLocations(ctx, &dispatchReq)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("couldn't resolve locations", zap.Error(err))
		return
	}

	// dispatch a driver request
	err = s.dispatcher.DispatchDriver(ctx, &dispatchReq, pickupLoc, dropoffLoc)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("couldn't trigger dispatch request", zap.Error(err))
		return
	}
	s.writeResponse(map[string]interface{}{}, w, r)
}

func (s *Server) notifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	sessionIDStr := r.Form.Get("sessionID")
	if sessionIDStr == "" {
		http.Error(w, "Missing required 'sessionID' parameter", http.StatusBadRequest)
		return
	}
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		http.Error(w, "Parameter 'sessionID' is not an integer", http.StatusBadRequest)
		return
	}

	cursorStr := r.Form.Get("cursor")
	if cursorStr == "" {
		http.Error(w, "Missing required 'cursor' parameter", http.StatusBadRequest)
		return
	}
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		http.Error(w, "Parameter 'cursor' is not an integer", http.StatusBadRequest)
		return
	}

	data, err := s.notification.List(ctx, uint(sessionID), cursor)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("request failed", zap.Error(err))
		return
	}

	s.writeResponse(data, w, r)
}

func (s *Server) writeResponse(response interface{}, w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(r.Context()).Error("cannot marshal response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
