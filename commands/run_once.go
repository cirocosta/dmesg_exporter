package commands

import (
	"context"
	"fmt"

	"github.com/cirocosta/dmesg_exporter/kmsg"
	"github.com/cirocosta/dmesg_exporter/reader"
)

type runOnce struct{}

func (c *runOnce) Execute(args []string) (err error) {
	r := reader.NewReader()

	ctx, cancel := context.WithCancel(context.Background())
	go blockAndCancelOnSignal(cancel)

	var (
		messages = make(chan *kmsg.Message, 1)
		errors   = r.Listen(ctx, messages)
	)

	for {
		select {
		case err = <-errors:
			return
		case message := <-messages:
			fmt.Println(message.Message)
		}
	}

	return
}
