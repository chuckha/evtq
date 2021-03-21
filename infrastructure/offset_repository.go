package infrastructure

import (
	"github.com/chuckha/evtq/core"
)

type OffsetRepository struct {
	// map [ connector name ] map [ event type ]
	offsets map[string]map[string]*core.Offset
}

func (o *OffsetRepository) GetOffsets(name string) []*core.Offset {
	out := []*core.Offset{}
	offsets, ok := o.offsets[name]
	if !ok {
		return out
	}
	for _, offset := range offsets {
		out = append(out, offset)
	}
	return out
}

func (o *OffsetRepository) AddOffset(name string, off *core.Offset) {
	o.offsets[name][off.EventType] = off
}
