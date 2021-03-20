package domain2

import (
	"io"
	"time"

	"github.com/pkg/errors"
)

type Decoding interface {
	NewDecoder(io.Reader) Decoder
}

type Decoder interface {
	Decode(interface{}) error
}

type Handler interface {
	Handle(evt *Event) error
}

type Consumer struct {
	Decoding
	Handler
}

func NewConsumer(dec Decoding, handler Handler) *Consumer {
	return &Consumer{
		Decoding: dec,
		Handler:  handler,
	}
}

func (c *Consumer) Consume(reader io.Reader) error {
	for {
		decoder := c.NewDecoder(reader)
		evt := &Event{}
		err := decoder.Decode(evt)
		if err == io.EOF {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if err != nil {
			return errors.WithStack(err)
		}
		if err := c.Handle(evt); err != nil {
			return err
		}
	}
}

type EventCountHandler struct {
	Count int
}

func (e *EventCountHandler) Handle(_ *Event) error {
	e.Count++
	return nil
}
