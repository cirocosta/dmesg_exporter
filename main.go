package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var command struct {
	TelemetryPath string   `long:"path" default:"/metrics" description:"path to serve metrics"`
	ListenAddress string   `long:"address" default:":9000" description:"address to listen for prometheus scraping"`
}

func main () {
	_, err := flags.Parse(command)
	if err != nil {
		os.Exit(1)
	}

	return
}
