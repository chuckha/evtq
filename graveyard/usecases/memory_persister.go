package usecases

import (
	"io"

	"github.com/chuckha/evtq/graveyard/domain"
)

type InMemPersist struct {
	events map[string][][]byte
}

func NewInMemPersist() *InMemPersist {
	return &InMemPersist{map[string][][]byte{}}
}

func (i *InMemPersist) PersistEvent(data []byte, eventType string) error {
	i.events[eventType] = append(i.events[eventType], data)
	return nil
}

type ReceiverManager struct {
	// event type to receiver
	receivers map[string][]domain.Receiver
	// receiver name to event list
	eventBuffer map[string][][]byte
}

func (r *ReceiverManager) NewEventNotification(data []byte, evtType string) {
	receivers := r.receivers[evtType]
	for _, receiver := range receivers {
		r.eventBuffer[receiver.Name()] = append(r.eventBuffer[receiver.Name()], data)
	}
}

func (r *ReceiverManager) GetEventsFor(receiverName string) [][]byte {
	return r.eventBuffer[receiverName]
}

// a thing that accepts events
//   when an event is accepted it is "persisted"
//   the new event is given to another process
// another process running
//   when it gets a new event it writes to a receiver. a receiver is an io.Writer

type AnotherProcess interface {
	EventSender([]byte, domain.Receiver) error
}

type AP struct{}

func (a *AP) EventSender(data []byte, receiever io.Writer) error {

}
