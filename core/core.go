package core

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

type EventStore interface {
	WriteEvent(e *Event)
	GetEventsFrom(offset int) []*Event
}

type Broker interface {
	AcceptEvent(e *Event)
}

type EventSender interface {
	WriteEvent(e *Event, writers []io.Writer) error
}

// EventSender
// writes to an io.Writer
//
