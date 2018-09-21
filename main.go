package main

import (
	"os"

	"github.com/cirocosta/dmesg_exporter/commands"
	"github.com/jessevdk/go-flags"
)

var err error

func main() {
	_, err = flags.Parse(&commands.DmesgExporter)
	if err != nil {
		os.Exit(1)
	}

	return
}
