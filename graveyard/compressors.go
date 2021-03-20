package graveyard

import (
	"encoding/json"

	"github.com/chuckha/evtq/graveyard/domain"
)

type DefaultCompressor struct{}

func NewDefaultCompressor() *DefaultCompressor {
	return &DefaultCompressor{}
}

func (c DefaultCompressor) Compress(evt *domain.Event) ([]byte, error) {
	o, err := json.Marshal(evt)
	if err != nil {
		return nil, err
	}
	return o, nil
}
