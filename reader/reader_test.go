package reader_test

import (
	"bytes"
	"context"
	"io"

	"github.com/cirocosta/dmesg_exporter/kmsg"
	"github.com/cirocosta/dmesg_exporter/reader"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	var (
		input        io.Reader
		errChan      chan error
		messagesChan chan *kmsg.Message
		r            reader.Reader
		err          error
	)

	JustBeforeEach(func() {
		r = reader.NewReader(input)
		messagesChan = make(chan *kmsg.Message, 1)

		errChan = r.Listen(
			context.Background(),
			messagesChan,
		)
	})

	Context("with malformed input", func() {
		BeforeEach(func() {
			buf := new(bytes.Buffer)
			_, err = buf.WriteString("something very wrong")

			Expect(err).NotTo(HaveOccurred())

			input = buf
		})

		It("sends failures over channel", func() {
			Eventually(errChan).Should(Receive(&err))
			Expect(err).To(HaveOccurred())
		})

		It("closes messages channel", func() {
			Eventually(messagesChan).Should(BeClosed())
		})
	})

	Context("with welformed input", func() {
		BeforeEach(func() {
			buf := new(bytes.Buffer)
			_, err = buf.WriteString("6,339,5140900,-;something")

			Expect(err).NotTo(HaveOccurred())

			input = buf
		})

		It("doesn't fail", func() {
			Consistently(errChan).ShouldNot(Receive(&err))
		})

		It("receives message", func() {
			var receivedMessage *kmsg.Message

			Eventually(messagesChan).Should(Receive(&receivedMessage))
			Expect(receivedMessage.Message).To(Equal("something"))
		})
	})
})
