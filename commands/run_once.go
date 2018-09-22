package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/cirocosta/dmesg_exporter/kmsg"
	"github.com/cirocosta/dmesg_exporter/reader"
)

const kmsgDevice = "/dev/kmsg"

type runOnce struct{
	Tail bool `long:"tail" short:"t" description:"whether the reader should seek to the end"`
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
