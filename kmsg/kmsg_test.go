package kmsg_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cirocosta/dmesg_exporter/kmsg"
)

var _ = Describe("Kmsg", func() {
	Describe("IsValidFacility", func() {
		It("validates if the facility is valid or not", func() {
			var (
				valid = []uint8{
					1, 2, 3, 4, 5, 6,
				}
				invalid = []uint8{
					30, 31, 230, 111,
				}
			)

			for _, facility := range valid {
				Expect(kmsg.IsValidFacility(facility)).To(BeTrue())
			}

			for _, facility := range invalid {
				Expect(kmsg.IsValidFacility(facility)).ToNot(BeTrue())
			}
		})
	})

	Describe("IsValidPriority", func() {
		It("validates if the priority is valid or not", func() {
			var (
				valid = []uint8{
					1, 2, 3, 4, 5, 6,
				}
				invalid = []uint8{
					30, 31, 230, 111,
				}
			)

			for _, priority := range valid {
				Expect(kmsg.IsValidPriority(priority)).To(BeTrue())
			}

			for _, priority := range invalid {
				Expect(kmsg.IsValidPriority(priority)).ToNot(BeTrue())
			}
		})
	})

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

			Context("not having a valid priority field", func() {
				BeforeEach(func() {
					input = "aa;aa"
				})

				It("fails", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
