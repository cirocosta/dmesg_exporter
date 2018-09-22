package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var DmesgExporter struct {
	RunOnce runOnce `command:"run-once"`
	Start   start   `command:"start"`
}

func blockAndCancelOnSignal(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	cancel()
}
