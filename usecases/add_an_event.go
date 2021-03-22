package usecases

import (
	"github.com/chuckha/evtq/core"
)

type EventWriter interface {
	WriteEvent(e *core.Event) error
}

type EventDistributor interface {
	DistributeEvent(e *core.Event)
}

type EventAdder struct {
	ew EventWriter
	ed EventDistributor
}

func (ea *EventAdder) AddEvent(evt *core.Event) error {
	if err := ea.ew.WriteEvent(evt); err != nil {
		return err
	}
	ea.ed.DistributeEvent(evt)
	return nil
}
