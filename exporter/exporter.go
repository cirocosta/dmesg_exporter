// Package exporter defines the implementation of dmesg_exporter's
// prometheus exporter internals, providing means for it to gather
// metrics from the registered collectors.
//
// ps.: The package is meant to be used by the main command only as it
//      doesn't provide any interface for generic loggers.
package exporter

import (
	"context"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ExporterConfig provides the configuration necessary to
// instantiate a new Exporter via `NewExporter`.
type Exporter struct {
	// ListenAddress is the address used by prometheus
	// to listen for scraping requests.
	//
	// Examples:
	// - :8080
	// - 127.0.0.2:1313
	ListenAddress string

	// TelemetryPath configures the path under which
	// the prometheus metrics are reported.
	//
	// For instance:
	// - /metrics
	// - /telemetry
	TelemetryPath string

	listener net.Listener
}

// Listen initiates the HTTP server using the configurations
// provided via ExporterConfig.
//
// This is a blocking method - make sure you either make use of
// goroutines to not block if needed.
func (e *Exporter) Listen(ctx context.Context) (err error) {
	http.Handle(e.TelemetryPath, promhttp.Handler())

	e.listener, err = net.Listen("tcp", e.ListenAddress)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to listen on address %s", e.ListenAddress)
		return
	}

	doneChan := make(chan error, 1)

	go func() {
		defer close(doneChan)

		err := http.Serve(e.listener, nil)
		if err != nil {
			doneChan <- errors.Wrapf(err,
				"failed listening on address %s",
				e.ListenAddress)
			return
		}
	}()

	select {
	case err = <-doneChan:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return
}

// Close gracefully closes the tcp listener associated with it.
func (e *Exporter) Close() (err error) {
	if e.listener == nil {
		return
	}

	err = e.listener.Close()
	return
}
