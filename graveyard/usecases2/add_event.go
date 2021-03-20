package usecases2

import (
	"github.com/chuckha/evtq/graveyard/domain2"
)

type ConsumerGetter interface {
	Get(evtType string) []*domain2.Consumer
}

type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

type EventAdder struct {
	ConsumerGetter
	Encoder
}

func (e *EventAdder) AddEvent(evt *domain2.Event) error {
	consumers := e.Get(evt.Type())
	encoded, err := e.Encode(evt)
	if err != nil {
		return err
	}
	for _, consumer := range consumers {
		if _, err := consumer.Write(encoded); err != nil {
			return err
		}
		// write to offset updater
	}
	return nil
}

/*
what is an offset?
after a write is complete, the offset is updated (perhaps simply inc by one)

when a new consumer joins, we inspect the latest known offset
	if the latest known offset is behind the newest event
	send every message between latest known and newest until the offset is at the latest


an offset it a marker that indicates if an event should be sent or not in an attempt to deliver at last once semantics

at least once...
could redeliver the entire queue on every reload

or could try to deliver from last known offset
last known offset...

an offset is a marker that identifies events in a list of events

since i suck at coding i'm using event number// line count as an offset
*/
