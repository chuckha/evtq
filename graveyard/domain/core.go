package domain

import (
	"time"

	"github.com/pkg/errors"
)

type Event struct {
	Created time.Time
	Type    string
	Data    []byte
}

func NewEvent(eventType string, data []byte) (*Event, error) {
	if eventType == "" {
		return nil, errors.New("Event type is required")
	}
	return &Event{
		Created: time.Now(),
		Type:    eventType,
		Data:    data,
	}, nil
}

// broker glues some behaviors together and exposes some fo them over one interface
type broker struct {
	compressor    Compressor
	offsetManager OffsetManager
	// ReceicerManager keeps track of active receivers and allows a fast lookup of all receivers by event type
	ReceiverManager

	// ReceiverMetadata keeps track of receivers and the last known offset for them.
}

// NewBroker is the only way to create a broker type.
func NewBroker(co Compressor, offsetM OffsetManager, receiverM ReceiverManager) (*broker, error) {
	if co == nil {
		return nil, errors.New("a DefaultCompressor is required")
	}
	if offsetM == nil {
		return nil, errors.New("an offset manager is required")
	}
	if receiverM == nil {
		return nil, errors.New("a receiver manager is requried")
	}
	b := &broker{
		compressor:      co,
		offsetManager:   offsetM,
		ReceiverManager: receiverM,
	}
	return b, nil
}

func (b *broker) Accept(evt *Event) error {
	data, offset, err := b.compressor.Compress(evt)
	if err != nil {
		return err
	}
	receivers := b.ReceiverManager.GetReceiver(evt.Type)
	for _, rcv := range receivers {
		if b.offsetManager.Get(ProcessorID(rcv.Name(), evt.Type)) >= offset {
			continue
		}
		if err := rcv.Receive(data); err != nil {
			return err
		}
		b.offsetManager.Set(ProcessorID(rcv.Name(), evt.Type), offset)
	}
	return nil
}

type PID string

func ProcessorID(rcvName, evtType string) PID {
	return PID(rcvName + "-" + evtType)
}
