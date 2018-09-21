package commands

import (
	"log"

	"github.com/cirocosta/dmesg_exporter/exporter"
)

type start struct {
	TelemetryPath string `long:"path" default:"/metrics" description:"path to serve metrics"`
	ListenAddress string `long:"address" default:":9000" description:"address to listen for prometheus scraping"`
}

func (c *start) Execute(args []string) (err error) {
	promExporter := exporter.Exporter{
		ListenAddress: DmesgExporter.Start.ListenAddress,
		TelemetryPath: DmesgExporter.Start.TelemetryPath,
		Collectors:    nil,
	}

	err = promExporter.Listen()
	if err != nil {
		log.Fatalln(err)
	}

	return
}
