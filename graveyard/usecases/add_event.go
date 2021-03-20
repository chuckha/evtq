package usecases

import (
	"github.com/chuckha/evtq/graveyard/domain"
)

type Notifier interface {
	NewEventNotification(data []byte, evtType string)
}

type Compressor interface {
	Compress(event *domain.Event) ([]byte, error)
}

type EventPersister interface {
	PersistEvent(data []byte, eventType string) error
}

type EventAdder struct {
	Compressor
	EventPersister
	Notifier
}

// TODO: add waitgroups to send the data to each receiver as quickly as possible
func (e *EventAdder) AddEvent(event *domain.Event) error {
	// - compress the event
	data, err := e.Compress(event)
	if err != nil {
		return err
	}
	// - persist the event
	if err := e.PersistEvent(data, event.Type); err != nil {
		return err
	}
	e.Notifier.NewEventNotification(data, event.Type)
	return nil
}

/*
log fan out

log1 [end filepointer, waiting at EOF, receiver pointers that are reading through until EOF then sleeping or reading and sending data wildly]
log2
log3
that all have the same "event type"

i write an event to the log

*/
