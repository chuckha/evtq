package domain2

type MemoryEventStore struct {
	store map[string][][]byte
}

func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		store: map[string][][]byte{},
	}
}

func (m *MemoryEventStore) getEventsFromOffset(etype string, offset int) ([][]byte, error) {
	return m.store[etype][offset:], nil
}

func (m *MemoryEventStore) getEventCount(etype string) int {
	return len(m.store[etype])
}

func (m *MemoryEventStore) write(etype string, d []byte) error {
	m.store[etype] = append(m.store[etype], d)
	return nil
}
