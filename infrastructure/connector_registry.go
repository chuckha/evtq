package infrastructure

import (
	"fmt"

	"github.com/chuckha/evtq/core"
)

type log interface {
	Debug(string)
}

type EventSender interface {
	SendEvents(*core.Event) error
}

type ConnectorsRegistry struct {
	connectors      map[string]core.Connector
	connectorsByEvt map[string]map[string]core.Connector
	log             log
}

type ConnectorsRegistryOption func(registry *ConnectorsRegistry)

func WithLog(log log) ConnectorsRegistryOption {
	return func(registry *ConnectorsRegistry) {
		registry.log = log
	}
}

func NewConnectorsRegistry(options ...ConnectorsRegistryOption) *ConnectorsRegistry {
	cr := &ConnectorsRegistry{
		connectors:      map[string]core.Connector{},
		connectorsByEvt: map[string]map[string]core.Connector{},
	}
	for _, o := range options {
		o(cr)
	}
	return cr
}

func (c *ConnectorsRegistry) RegisterConnector(connector core.Connector) {
	c.connectors[connector.GetName()] = connector
	for _, offset := range connector.GetOffsets() {
		if _, ok := c.connectorsByEvt[offset.EventType]; !ok {
			c.connectorsByEvt[offset.EventType] = make(map[string]core.Connector)
		}
		c.connectorsByEvt[offset.EventType][connector.GetName()] = connector
	}
}

// DistributeEvent sends an event to every interested connector
// It does not return an error ever as that may impede event flow.
// Instead, remove the problematic connectors
func (c *ConnectorsRegistry) DistributeEvent(e *core.Event) {
	toRemove := []core.Connector{}
	for _, connector := range c.connectorsByEvt[e.EventType] {
		if err := connector.SendEvents(e); err != nil {
			c.log.Debug(fmt.Sprintf("encountered an error with connector %s\n: %v", connector.GetName(), err))
			toRemove = append(toRemove, connector)
		}
	}
	for _, remove := range toRemove {
		c.removeConnector(remove)
	}
}

func (c *ConnectorsRegistry) removeConnector(remove core.Connector) {
	delete(c.connectors, remove.GetName())
	for _, connectorMap := range c.connectorsByEvt {
		delete(connectorMap, remove.GetName())
	}
}
