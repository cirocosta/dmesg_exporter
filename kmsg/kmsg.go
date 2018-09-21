// Package kmsg provides a minimal interface for dealing with
// kmsg messages extracted from `/dev/kmesg`.
package kmsg

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Level uint8

const (
	LevelEmerg Level = iota
	LevelAlert
	LevelCrit
	LevelErr
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug
)

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
)

type Message struct {
	Level          Level
	Facility       Facility
	SequenceNumber int64
	Timestamp      time.Time
	Message        string
	Metadata       map[string]string
}

// Parse takes care of parsing a `kmsg` message acording to the kernel
// documentation at https://www.kernel.org/doc/Documentation/ABI/testing/dev-kmsg.
//
// REGULAR MESSAGE:
//
//                  INFO												       MSG
//     .------------------------------------------. .------.
//    |																		         |			  |
//    |																		         |			  |
//    |	int	    int      int       char, <ignore>  | string |
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
	if len(splittedMessage) == 0 {
		err = errors.Errorf("message field not present")
		return
	}

	return
}
