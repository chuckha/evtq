package usecases2

import (
	"github.com/pkg/errors"

	"github.com/chuckha/evtq/graveyard/domain2"
)

type ConsumerSaver interface {
	Save(c *domain2.Consumer) error
}

type WriterAdder struct {
	domain2.ConsumerBuilder
	ConsumerSaver
}

func (a *WriterAdder) AddConsumer(wbi *domain2.ConsumerBuilderInfo) error {
	consumer, err := a.Build(wbi)
	if err != nil {
		return err
	}
	return errors.WithStack(a.Save(consumer))
	// TODO: send all messages not received (from offset to latest)
	// lookup known offset
	// 	if there is no known offset, set the last known offset to 0
	// using the last known offset, create a reader at offset
	// begin reading and sending messages until the offset is == latest

	// create the thing that perpetually reads new events and writes to writer
}

// Consumer layer
// WriteEventToConsumer(evt, cons)

// public applicaiton layer
// AcceptEvent(evt) -> event layer
// (writing event to event layer)
// NewConsumer(consumerinfo) -> event layer
// (creating a consumer, writing any unreceived events, call bind)

// private application layer
// BindEventAndConsumerLayer ->
// 	creating a forever process that reads from the event layer and writes to the consumer layer

// Event layer
// PerpetuallyReadEventsOfType(etype) io.Reader
// 		how? if it's files, we read the file, don't ever return eof, sleep repeat
//      if it's an in memory buffer, we read the buffer forever
// WriteEvent(evt)
// ReadEventsFromOffset(etype, offset) io.Reader
