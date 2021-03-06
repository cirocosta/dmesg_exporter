package kmsg_test

import (
	"time"

	"github.com/cirocosta/dmesg_exporter/kmsg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kmsg", func() {
	Describe("DecodeFlag", func() {
		It("decodes the flag accordingly", func() {
			var (
				actual    kmsg.Flag
				testCases = []struct {
					input    byte
					expected kmsg.Flag
				}{
					{
						input:    'c',
						expected: kmsg.FlagFragment,
					},
					{
						input:    '-',
						expected: kmsg.FlagDefault,
					},
					{
						input:    '+',
						expected: kmsg.FlagFragmentFollowing,
					},
					{
						input:    'z',
						expected: kmsg.FlagUnknown,
					},
				}
			)

			for _, tc := range testCases {
				actual = kmsg.DecodeFlag(tc.input)
				Expect(actual).To(Equal(tc.expected))
			}
		})
	})

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

	Describe("DecodePrefix", func() {
		var (
			prefix   uint8
			priority kmsg.Priority
			facility kmsg.Facility
		)

		JustBeforeEach(func() {
			priority, facility = kmsg.DecodePrefix(prefix)
		})

		Context("with unknown facility encoded in the prefix", func() {
			BeforeEach(func() {
				prefix = (1 << 6)
			})

			It("returns unknown facility", func() {
				Expect(facility).To(Equal(kmsg.FacilityUnknown))
				Expect(priority).To(Equal(kmsg.PriorityEmerg))
			})
		})

		Context("with known facility and priority", func() {
			BeforeEach(func() {
				prefix = 1
				prefix = prefix | 1<<3
			})

			It("returns the proper facility", func() {
				Expect(facility).To(Equal(kmsg.FacilityUser))
			})

			It("returns the proper priorty", func() {
				Expect(priority).To(Equal(kmsg.PriorityAlert))
			})
		})
	})

	Describe("Parse", func() {
		var (
			input   string
			message *kmsg.Message
			err     error
		)

		JustBeforeEach(func() {
			message, err = kmsg.Parse(input)
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

			Context("having message and info fields", func() {
				Context("not having enough fields in info section", func() {
					BeforeEach(func() {
						input = "aaaa;aaa"
					})

					It("fails", func() {
						Expect(err).To(HaveOccurred())
					})
				})

				Context("having enough fields in info section", func() {
					Context("with proper fields", func() {
						var expectedMessage *kmsg.Message

						BeforeEach(func() {
							input = "1,1,100,-,anything;message\n\n"
							expectedMessage = &kmsg.Message{
								Facility:       kmsg.FacilityKern,
								Flag:           kmsg.FlagDefault,
								Message:        "message",
								Priority:       kmsg.PriorityAlert,
								SequenceNumber: 1,
								Timestamp:      time.Unix(100/int64(time.Millisecond), 0),
							}
						})

						It("properly parses the fields", func() {
							Expect(message).To(BeEquivalentTo(expectedMessage))
						})
					})

					Context("with malformed prefix", func() {
						Context("being a string", func() {
							BeforeEach(func() {
								input = "a,b,c,d,e;message"
							})

							It("fails", func() {
								Expect(err).To(HaveOccurred())
							})
						})

						Context("being a int bigger than uint8", func() {
							BeforeEach(func() {
								input = "999999999,b,c,d,e;message"
							})

							It("fails", func() {
								Expect(err).To(HaveOccurred())
							})
						})

						Context("being a negative integer", func() {
							BeforeEach(func() {
								input = "-9999,b,c,d,e;message"
							})

							It("fails", func() {
								Expect(err).To(HaveOccurred())
							})
						})
					})

					Context("with malformed sequence", func() {
						BeforeEach(func() {
							input = "1,b,c,d,e;message"
						})

						It("fails", func() {
							Expect(err).To(HaveOccurred())
						})
					})

					Context("with malformed timestamp", func() {
						BeforeEach(func() {
							input = "1,1,c,d,e;message"
						})

						It("fails", func() {
							Expect(err).To(HaveOccurred())
						})
					})

					Context("with malformed flag", func() {
						BeforeEach(func() {
							input = "1,1,1,,e;message"
						})

						It("fails", func() {
							Expect(err).To(HaveOccurred())
						})
					})
				})
			})

		})
	})
})
