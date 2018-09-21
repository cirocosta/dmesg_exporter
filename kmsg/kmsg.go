// Package kmsg provides a minimal interface for dealing with
// kmsg messages extracted from `/dev/kmesg`.
package kmsg

import (
	"errors"
	"time"
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

var (
	ErrMessageMetadata    = errors.New("metadata message continuation not supported")
	ErrMessageInBadFormat = errors.New("unsuported format")
)

func Parse(rawMsg string) (m *Message, err error) {
	return
}
