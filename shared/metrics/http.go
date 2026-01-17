package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/usc-platform/shared/logging"
)

// StartHTTPServer starts a Prometheus metrics HTTP server on the given address and path.
// It returns a shutdown function to gracefully stop the server.
func StartHTTPServer(addr string, path string, logger *logging.Logger) (func(ctx context.Context) error, error) {
	if path == "" {
		path = "/metrics"
	}

	mux := http.NewServeMux()
	mux.Handle(path, promhttp.Handler())

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		logger.Info("Starting metrics HTTP server",
			logging.String("address", addr),
			logging.String("path", path))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics HTTP server failed",
				logging.String("address", addr),
				logging.Error(fmt.Errorf("listen and serve: %w", err)))
		}
	}()

	return func(ctx context.Context) error {
		logger.Info("Shutting down metrics HTTP server",
			logging.String("address", addr))
		return server.Shutdown(ctx)
	}, nil
}
