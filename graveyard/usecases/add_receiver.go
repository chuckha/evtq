package usecases

import (
	"github.com/chuckha/evtq/graveyard/domain"
)

type Metadata struct {
	LastKnownOffsets map[string]int
}

type EventSender interface {
	SendEvents(eventType string, from int) error
}

type ReceiverManager interface {
	Lookup(rcv domain.Receiver) (Metadata, error)
}

type ReceiverAdder struct {
	ReceiverManager
	EventSender
}

func (r *ReceiverAdder) AddReceiver(rcv domain.Receiver) error {
	for _, eventType := range rcv.EventTypes() {
		latestOffset := r.LookupLatestOffset(rcv.Name, eventType)
		missingEvents := r.StartReadingEvents(type, latestOffset)
		r.SendEvents(missingEvents)
	}
	return nil
}

// for each event type
// 1. look up its latest offset
// 2. events := getEvents(type, start) // get events returns a list of events from start to finish
// 3. send(events, r)
// 4. update offset for (receiver,eventtype)

// when a receiver attaches
// we look at each event it cares about
// we create a reader
// we gather events from the reader (current offset is required)
// we send events to the receiver
