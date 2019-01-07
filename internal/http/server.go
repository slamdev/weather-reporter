package http

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Router struct {
	Method  string
	Path    string
	Queries []string
	Handler http.Handler
}

type HttpServer interface {
	Start() error
	Stop() error
}

func NewHttpServer(serverPort int, routers ...Router) HttpServer {
	return &httpServer{
		server: http.Server{
			Addr:    ":" + strconv.Itoa(serverPort),
			Handler: buildRootHandler(routers...),
		},
	}
}

type httpServer struct {
	server http.Server
}

func (s *httpServer) Start() error {
	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		// do not consider server closing as error
		return nil
	}
	return errors.Wrap(s.server.ListenAndServe(), "failed to start server")
}

func (s *httpServer) Stop() error {
	return errors.Wrap(s.server.Shutdown(context.Background()), "failed to stop server")
}

func buildRootHandler(routers ...Router) http.Handler {
	rootRouter := mux.NewRouter()
	rootRouter.HandleFunc("/health", handleHealthRequest).Methods("GET")
	rootRouter.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})).Methods("GET")
	for _, router := range routers {
		rootRouter.NewRoute().Methods(router.Method).Path(router.Path).Queries(router.Queries...).
			Handler(instrumentHandler(router.Handler))
	}
	rootRouter.Use(logRequestsHandler, handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	return rootRouter
}

func handleHealthRequest(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(writer, `{"status": "SERVING"}`)
	if err != nil {
		log.Errorln(err)
	}
}

func instrumentHandler(handler http.Handler) http.Handler {
	return promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handler)
}

func logRequestsHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		for _, skipPath := range []string{"/health", "/metrics"} {
			if request.URL.Path == skipPath {
				handler.ServeHTTP(writer, request)
				return
			}
		}
		startTime := time.Now()
		recorder := statusCodeRecorder{writer, http.StatusOK}
		handler.ServeHTTP(&recorder, request)
		log.WithFields(log.Fields{
			"elapsedTime": time.Since(startTime),
			"requestIP":   getRequestIPAddress(request),
			"requestPath": request.URL.Path,
			"statusCode":  recorder.statusCode,
		}).Info("request")
	})
}

type statusCodeRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusCodeRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func getRequestIPAddress(request *http.Request) string {
	address := request.Header.Get("X-Forwarded-For")
	if len(address) == 0 {
		address = request.Header.Get("X-Real-IP")
	}
	if len(address) == 0 {
		var err error
		address, _, err = net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			log.WithField("error", fmt.Sprintf("%+v", err)).
				WithField("addr", request.RemoteAddr).
				Error("failed to split remote address")
			return ""
		}
	}
	return address
}
