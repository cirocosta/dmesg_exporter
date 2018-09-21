// Package kmsg provides a minimal interface for dealing with
// kmsg messages extracted from `/dev/kmesg`.
package kmsg

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Priority uint8

const (
	PriorityEmerg Priority = iota
	PriorityAlert
	PriorityCrit
	PriorityErr
	PriorityWarning
	PriorityNotice
	PriorityInfo
	PriorityDebug
	priorityEnd
)

func IsValidPriority (priority uint8) (isValid bool) {
	isValid = (priority < uint8(priorityEnd))
	return
}

type Facility uint8

const (
	FacilityKern Facility = iota
	FacilityUser
	FacilityMail
	FacilityDaemon
	FacilityAuth
	FacilitySyslog
	FacilityLpr
	FacilityNews
	facilityEnd
)

func IsValidFacility (facility uint8) (isValid bool) {
	isValid = (facility < uint8(facilityEnd))
	return
}

type Message struct {
	Priority          Priority
	Facility       Facility
	SequenceNumber int64
	Timestamp      time.Time
	Message        string
	Metadata       map[string]string
}

// DecodePrefix extracts both priority and facility from a given
// syslog(2) encoded prefix.
//
//	    facility    priority
//      .--------------.-----.
//      |              |     |
//	7  6  5  4  3  2  1  0    bits
//
//
func DecodePrefix(prefix uint8) (priority Priority, facility Facility, err error) {
	priority = Priority((prefix << (8-3)) >> 8-3)
	facility = Facility(prefix >> 3)

	return
}

// Parse takes care of parsing a `kmsg` message acording to the kernel
// documentation at https://www.kernel.org/doc/Documentation/ABI/testing/dev-kmsg.
//
// REGULAR MESSAGE:
//
//                  INFO		              MSG
//     .------------------------------------------. .------.
//    |                                            |        |
//    |	int	int      int      char, <ignore>   | string |
//    priority, seq, timestamp_us,flag[,..........];<message>
//
//
// CONTINUATION:
//
//	    | key | value |
//	/x7F<THIS>=<THATTT>
//
func Parse(rawMsg string) (m *Message, err error) {
	if rawMsg == "" {
		err = errors.Errorf("msg must not be empty")
		return
	}

	splittedMessage := strings.SplitN(rawMsg, ";", 2)
	if len(splittedMessage) < 2 {
		err = errors.Errorf("message field not present")
		return
	}

	m = new(Message)

	infoSection := splittedMessage[0]
	m.Message = splittedMessage[1]

	splittedInfoSection := strings.SplitN(infoSection, ",", 5)
	if len(splittedInfoSection) < 4 {
		err = errors.Errorf("info section with not enought fields")
		return
	}

	_, err = strconv.Atoi(splittedInfoSection[0])
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't convert priority to int")
		return
	}

	// CC: parse info section

	return
}
