package domain2

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

// Event and fields must be public to use encoding packages
type Event struct {
	Created   time.Time
	EventType string
	Data      []byte
}

func (e *Event) String() string {
	return fmt.Sprintf("[%s]: %s", e.EventType, e.Data)
}

func NewEvent(eventType string, data []byte) (*Event, error) {
	if eventType == "" {
		return nil, errors.New("Event type is required")
	}
	return &Event{
		Created:   time.Now(),
		EventType: eventType,
		Data:      data,
	}, nil
}

func (e *Event) Type() string {
	return e.EventType
}

type EventStore interface {
	Write(evt *Event) error
	EventsFromOffset(etype string, offset int) io.Reader
	EventListener(etype string) io.Reader
}

// PerpetuallyReadEventsOfType(EventType) io.Reader
// 		how? if it's files, we read the file, don't ever return eof, sleep repeat
//      if it's an in memory buffer, we read the buffer forever
// WriteEvent(evt)
// ReadEventsFromOffset(EventType, offset) io.Reader
