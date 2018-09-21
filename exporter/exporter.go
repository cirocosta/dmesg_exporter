// Package exporter defines the implementation of dmesg_exporter's
// prometheus exporter internals, providing means for it to gather
// metrics from the registered collectors.
//
// ps.: The package is meant to be used by the main command only as it
//      doesn't provide any interface for generic loggers.
package exporter

import (
	"net"
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
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

	// Collectors holds a listen of collectors that are
	// meant to send metrics.
	Collectors []prometheus.Collector

	once     sync.Once
	listener net.Listener
}

func (e *Exporter) init() (err error) {
	for _, collector := range e.Collectors {
		err = prometheus.Register(collector)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to register collector")
			return
		}
	}

	return
}

// Listen initiates the HTTP server using the configurations
// provided via ExporterConfig.
//
// This is a blocking method - make sure you either make use of
// goroutines to not block if needed.
func (e *Exporter) Listen() (err error) {
	e.once.Do(func() {
		err = e.init()
	})
	if err != nil {
		err = errors.Wrapf(err,
			"failed to initialize exporter")
		return
	}

	http.Handle(e.TelemetryPath, promhttp.Handler())

	e.listener, err = net.Listen("tcp", e.ListenAddress)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to listen on address %s", e.ListenAddress)
		return
	}

	err = http.Serve(e.listener, nil)
	if err != nil {
		err = errors.Wrapf(err,
			"failed listening on address %s",
			e.ListenAddress)
		return
	}

	return
}

// Stop closes the tcp listener (if exists).
func (e *Exporter) Stop() (err error) {
	if e.listener == nil {
		return
	}

	err = e.listener.Close()
	return
}
