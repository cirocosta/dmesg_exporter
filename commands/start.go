package commands

import (
	"context"
	"os"

	"github.com/cirocosta/dmesg_exporter/exporter"
	"github.com/cirocosta/dmesg_exporter/kmsg"
	"github.com/cirocosta/dmesg_exporter/reader"
	"github.com/prometheus/client_golang/prometheus"
)

type start struct {
	TelemetryPath string `long:"path" default:"/" description:"path to serve metrics" env:"PATH"`
	ListenAddress string `long:"address" default:":9000" description:"address to listen for prometheus scraping" env:"ADDRESS"`
	Tail          bool   `long:"tail" description:"seek to the end when starting" env:"TAIL"`
}

func (c *start) Execute(args []string) (err error) {
	file, err := os.Open(kmsgDevice)
	if err != nil {
		return
	}
	defer file.Close()

	if c.Tail {
		_, err = file.Seek(0, os.SEEK_END)
		if err != nil {
			return
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go blockAndCancelOnSignal(cancel)

	var (
		r                  = reader.NewReader(file)
		exporterErrorsChan = make(chan error, 1)
		messages           = make(chan *kmsg.Message, 1)
		promExporter       = exporter.Exporter{
			ListenAddress: DmesgExporter.Start.ListenAddress,
			TelemetryPath: DmesgExporter.Start.TelemetryPath,
		}
		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dmesg_logs",
				Help: "total number of logs received for a given facility and priority",
			},
			[]string{"priority", "facility"},
		)
	)

	kmsgErrorsChan := r.Listen(ctx, messages)
	prometheus.MustRegister(counter)

	defer promExporter.Close()

	go func() {
		defer close(exporterErrorsChan)

		err = promExporter.Listen(ctx)
		if err != nil {
			exporterErrorsChan <- err
		}
	}()

	for {
		select {
		case err = <-kmsgErrorsChan:
			return
		case err = <-exporterErrorsChan:
			return
		case message := <-messages:
			if message == nil {
				return
			}

			counter.With(prometheus.Labels{
				"priority": message.Priority.String(),
				"facility": message.Facility.String(),
			}).Inc()
		}
	}

	return
}
