package http

import (
	"errors"
	"fmt"
	"github.com/intezya/pkglib/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

const (
	writeTimeout = 5 * time.Second
	readTimeout  = 5 * time.Second
	idleTimeout  = 10 * time.Second
)

func SetupMetricsServer(port int) {
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		server := &http.Server{
			//nolint:exhaustruct
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      mux,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		}

		logger.Log.Infof("Starting metrics server on port %d", port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Warnf("Metrics server error: %v", err)
		}
	}()
}
