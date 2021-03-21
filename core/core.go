package core

import (
	"fmt"
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
