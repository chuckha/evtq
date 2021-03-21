package usecases

import (
	"github.com/chuckha/evtq/core"
)

type ConsumerBuilder interface {
	Build(*core.ConnectorBuilderInfo) (*core.Connector, error)
}

// ConnectorRepository must live in memory due to I/O and possible network connections or shared buffer
type ConnectorRepository interface {
	RegisterConnector(connector *core.Connector)
	DistributeEvent(e *core.Event) error
}

type EventStore interface {
	WriteEvent(e *core.Event) error
	GetEventsFrom(eventType string, offset int) ([]*core.Event, error) // TODO: change interface to accept *Offset
}

type ConsumerConnector struct {
	ConnectorRepository
	EventStore
}

func (c *ConsumerConnector) AddConnector(connector *core.Connector) error {
	// connect the connector
	c.ConnectorRepository.RegisterConnector(connector)

	for _, offset := range connector.Offsets {
		events, err := c.EventStore.GetEventsFrom(offset.EventType, offset.LastKnownOffset)
		if err != nil {
			return err
		}
		if err := connector.SendEvents(events...); err != nil {
			return err
		}
	}
	// update offset database
	return nil
}
