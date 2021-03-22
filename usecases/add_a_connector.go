package usecases

import (
	"io"

	"github.com/chuckha/evtq/core"
)

type ConsumerBuilder interface {
	Build(*core.ConnectorBuilderInfo) (core.Connector, error)
}

// ConnectorRegistry must live in memory due to I/O and possible network connections or shared buffer
type ConnectorRegistry interface {
	RegisterConnector(connector core.Connector)
}

type EventGetter interface {
	GetEventsFrom(eventType string, offset int) ([]*core.Event, error) // TODO: change interface to accept *Offset
}

type OffsetRepository interface {
	UpdateOffset(name string, offset *core.Offset)
}

func NewConnectorAdder(connectors ConnectorRegistry, events EventGetter, offsets OffsetRepository) *ConnectorAdder {
	return &ConnectorAdder{
		connectors: connectors,
		events:     events,
		offsets:    offsets,
	}
}

type ConnectorAdder struct {
	connectors ConnectorRegistry
	events     EventGetter
	offsets    OffsetRepository
}

// TODO: Change this signature to return an io.Reader for local ones and a message for remote ones telling them
// to be listening on a localhost port for events.
func (c *ConnectorAdder) AddConnector(connector core.Connector) (io.Reader, error) {
	c.connectors.RegisterConnector(connector)

	for _, offset := range connector.GetOffsets() {
		events, err := c.events.GetEventsFrom(offset.EventType, offset.LastKnownOffset)
		if err != nil {
			return nil, err
		}
		if err := connector.SendEvents(events...); err != nil {
			return nil, err
		}
		newOffset := core.NewOffset(offset.EventType, offset.LastKnownOffset+len(events))
		c.offsets.UpdateOffset(connector.GetName(), newOffset)
	}
	return connector.GetReader(), nil
}
