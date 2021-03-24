package evtq

import (
	"io"

	"github.com/chuckha/evtq/core"
	"github.com/chuckha/evtq/usecases"
)

type Broker struct {
	adapter
	useCases
}

func NewBroker(adapter adapter, useCases useCases) *Broker {
	return &Broker{
		adapter:  adapter,
		useCases: useCases,
	}
}

// AddConnector returns an io.Reader, but that's only to be used in the case of a
// local connector. Don't use it if you've asked for a TCP connector. Listen to the
// socket instead.
func (b *Broker) AddConnector(info *core.ConnectorBuilderInfo) (io.Reader, error) {
	connector, err := b.adapter.AddConnectorFromInfo(info)
	if err != nil {
		return nil, err
	}
	return b.useCases.AddConnector(connector)
}

func (b *Broker) AddEvent(eventType string, data []byte) error {
	evt, err := b.adapter.AddEvent(eventType, data)
	if err != nil {
		return err
	}
	return b.useCases.AddEvent(evt)
}

type adapter interface {
	AddConnectorFromInfo(info *core.ConnectorBuilderInfo) (core.Connector, error)
	AddEvent(eventType string, data []byte) (*core.Event, error)
}

type useCases interface {
	AddConnector(connector core.Connector) (io.Reader, error)
	AddEvent(evt *core.Event) error
}

type ConnectorBuilder interface {
	BuildConnector(info *core.ConnectorBuilderInfo) (core.Connector, error)
}

type BrokerAdapter struct {
	ConnectorBuilder
}

func (b *BrokerAdapter) AddConnectorFromInfo(info *core.ConnectorBuilderInfo) (core.Connector, error) {
	return b.BuildConnector(info)
}

func (b *BrokerAdapter) AddEvent(eventType string, data []byte) (*core.Event, error) {
	return core.NewEvent(eventType, data)
}

type BrokerUseCases struct {
	*usecases.ConnectorAdder
	*usecases.EventAdder
}
