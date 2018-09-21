package kmsg_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cirocosta/dmesg_exporter/kmsg"
)

var _ = Describe("Kmsg", func() {
	Describe("Parse", func() {
		var (
			input string
			// message *kmsg.Message
			err error
		)

		JustBeforeEach(func() {
			_, err = kmsg.Parse(input)
		})

		Context("with empty string", func() {
			It("fails", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with malformed string", func() {
			Context("not having message field", func() {
				BeforeEach(func() {
					input = "aaaaaaaa"
				})

				It("fails", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("having message field", func() {
				Context("not having enough fields in info section", func() {
					BeforeEach(func() {
						input = "aaaa;aaa"
					})

					It("fails", func() {
						Expect(err).To(HaveOccurred())
					})
				})
			})

		})
	})
})
