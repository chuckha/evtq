package core

import (
	"io"
)

type ConnectorInfo interface {
	Info()
}

type EncodingType string

type ConnectorBuilderInfo struct {
	Name         string
	EventTypes   []string
	EncodingType EncodingType
	Info         ConnectorInfo
}

type Offset struct {
	EventType       string
	LastKnownOffset int
}

type EncDec interface {
	Encoding
	Decoding
}

type Encoding interface {
	NewEncoder(io.Writer) Encoder
}

type Decoding interface {
	NewDecoder(io.Reader) Decoder
}

type Encoder interface {
	Encode(interface{}) error
}
type Decoder interface {
	Decode(interface{}) error
}

type Connector struct {
	Name    string
	IO      io.ReadWriter
	EncDec  EncDec
	Offsets []*Offset
}

func (c *Connector) GetName() string {
	return c.Name
}

func (c *Connector) GetOffsets() []*Offset {
	return c.Offsets
}

func (c *Connector) SendEvents(events ...*Event) error {
	for _, event := range events {
		err := c.EncDec.NewEncoder(c.IO).Encode(event)
		if err != nil {
			return err
		}
	}
	return nil
}
