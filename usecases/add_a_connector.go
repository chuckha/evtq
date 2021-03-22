package usecases

import (
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

func (c *ConnectorAdder) AddConnector(connector core.Connector) error {
	c.connectors.RegisterConnector(connector)

	for _, offset := range connector.GetOffsets() {
		events, err := c.events.GetEventsFrom(offset.EventType, offset.LastKnownOffset)
		if err != nil {
			return err
		}
		if err := connector.SendEvents(events...); err != nil {
			return err
		}
		newOffset := core.NewOffset(offset.EventType, offset.LastKnownOffset+len(events))
		c.offsets.UpdateOffset(connector.GetName(), newOffset)
	}
	return nil
}
