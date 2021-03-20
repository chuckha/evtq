package domain2

import (
	"fmt"
	"io"
	"os"
	"time"
)

type EncDec interface {
	Encoding
	Decoding
}

type Encoding interface {
	NewEncoder(io.Writer) Encoder
}

type Encoder interface {
	Encode(interface{}) error
}

type BackendStore interface {
	getEventsFromOffset(etype string, offset int) ([]*Event, error)
	getEventCount(etype string) int
	write(evt *Event) error
}

type Store struct {
	BackendStore
}

func NewStore(bs BackendStore) *Store {
	return &Store{
		BackendStore: bs,
	}
}

func (m *Store) Write(evt *Event) error {
	return m.write(evt)
}

func (m *Store) EventsFromOffset(etype string, offset int) ([]*Event, error) {
	return m.getEventsFromOffset(etype, offset)
}

func (m *Store) EventListener(etype string) chan *Event {
	eventc := make(chan *Event)
	offset := m.getEventCount(etype)
	go func(offset int) {
		ticker := time.Tick(100 * time.Millisecond)
		for {
			select {
			// TODO: case <-closec:
			case <-ticker:
				if m.getEventCount(etype) <= offset {
					continue
				}
				unreadEvents, err := m.getEventsFromOffset(etype, offset)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "failed to get events from a specific offset: %v\n", err)
					continue
				}
				for _, event := range unreadEvents {
					eventc <- event
				}
				offset = m.getEventCount(etype)
			}
		}
	}(offset)
	return eventc
}
