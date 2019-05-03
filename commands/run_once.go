package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/cirocosta/dmesg_exporter/kmsg"
	"github.com/cirocosta/dmesg_exporter/reader"
)

const kmsgDevice = "/dev/kmsg"

type runOnce struct {
	Tail bool `long:"no-tail" description:"seek to the end when starting" env:"TAIL"`
}

func (c *runOnce) Execute(args []string) (err error) {
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
		messages = make(chan *kmsg.Message, 1)
		r        = reader.NewReader(file)
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
