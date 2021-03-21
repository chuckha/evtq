package infrastructure

import (
	"github.com/chuckha/evtq/core"
)

type MemoryStore struct {
	events map[string][]*core.Event
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{map[string][]*core.Event{}}
}
func (m *MemoryStore) WriteEvent(e *core.Event) error {
	if _, ok := m.events[e.EventType]; !ok {
		m.events[e.EventType] = make([]*core.Event, 0)
	}
	m.events[e.EventType] = append(m.events[e.EventType], e)
	return nil
}

func (m *MemoryStore) GetEventsFrom(eventType string, eventNumber int) ([]*core.Event, error) {
	if len(m.events[eventType]) < eventNumber {
		return nil, nil
	}
	return m.events[eventType][eventNumber:], nil
}

// dump to disk?
// keep track of file, event #s and offsets
// what do we keep in memory?
// persist every 100 events?
// purge memory to disk and start over
// events from would have to
