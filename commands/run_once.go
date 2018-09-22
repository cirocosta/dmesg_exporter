package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/cirocosta/dmesg_exporter/kmsg"
	"github.com/cirocosta/dmesg_exporter/reader"
)

const kmsgDevice = "/dev/kmsg"

type runOnce struct{}

func (c *runOnce) Execute(args []string) (err error) {
	file, err := os.Open(kmsgDevice)
	if err != nil {
		return
	}
	defer file.Close()

	r := reader.NewReader(file)

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
			if message == nil {
				return
			}

			fmt.Println(message.Message)
		}
	}

	return
}
