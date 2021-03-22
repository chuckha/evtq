package core

import (
	"io"

	"github.com/pkg/errors"
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

func NewOffset(eventType string, offset int) *Offset {
	return &Offset{
		EventType:       eventType,
		LastKnownOffset: offset,
	}
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

type Connector interface {
	GetName() string
	GetOffsets() map[string]*Offset
	SendEvents(e ...*Event) error
	GetReader() io.Reader
}

func NewConnector(name string, inout io.ReadWriter, encdec EncDec, offsets map[string]*Offset) (*connector, error) {
	if name == "" {
		return nil, errors.New("connectors require names")
	}
	if len(offsets) == 0 {
		return nil, errors.New("connectors must watch at least one event type")
	}
	return &connector{
		Name:    name,
		IO:      inout,
		EncDec:  encdec,
		Offsets: offsets,
	}, nil
}

type connector struct {
	Name    string
	IO      io.ReadWriter
	EncDec  EncDec
	Offsets map[string]*Offset
}

func (c *connector) GetName() string {
	return c.Name
}

func (c *connector) GetOffsets() map[string]*Offset {
	return c.Offsets
}

func (c *connector) GetReader() io.Reader {
	return c.IO
}

func (c *connector) SendEvents(events ...*Event) error {
	for _, event := range events {
		err := c.EncDec.NewEncoder(c.IO).Encode(event)
		if err != nil {
			return err
		}
	}
	return nil
}
