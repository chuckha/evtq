package graveyard

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/chuckha/evtq/graveyard/domain"
)

type EventWriter interface {
	WriteEvent(data []byte, evtType string) error
}

type PersistenceCompressor struct {
	EventWriter
	domain.MaxOffsetTracker
}

func NewPersistenceCompressor(directory string) (*PersistenceCompressor, error) {
	ew, err := NewEventFileWriter(directory)
	if err != nil {
		return nil, err
	}
	return &PersistenceCompressor{
		EventWriter:      ew,
		MaxOffsetTracker: domain.NewDefaultMaxOffsetTracker(),
	}, nil
}

// TODO: figure out the compression swapping in a bit
func (pc *PersistenceCompressor) Compress(evt *domain.Event) ([]byte, int, error) {
	b, err := json.Marshal(evt)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	if err := pc.WriteEvent(b, evt.Type); err != nil {
		return nil, 0, err
	}
	return b, pc.Update(evt.Type, len(b)), nil
}
