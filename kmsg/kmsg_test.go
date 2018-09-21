package kmsg_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cirocosta/dmesg_exporter/kmsg"
)

var _ = Describe("Kmsg", func() {
	Describe("Parse", func () {
		var (
			input string
			// message *kmsg.Message
			err error
		)

		JustBeforeEach(func () {
			_, err = kmsg.Parse(input)
		})

		Context("with empty string", func () {
			It("fails", func () {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with malformed string", func () {
			BeforeEach(func () {
				input = "$$$$$something very malformed"
			})

			It("fails", func () {
				Expect(err).To(HaveOccurred())
			})

			It("describes the error", func () {
				Expect(err).To(Equal(kmsg.ErrMessageInBadFormat))
			})
		})

		Context("with metadata string", func () {
			BeforeEach(func () {
				input = " FOO=BAR"
			})

			It("fails", func () {
				Expect(err).To(HaveOccurred())
			})

			It("describes the error", func () {
				Expect(err).To(Equal(kmsg.ErrMessageMetadata))
			})
		})
	})
})
