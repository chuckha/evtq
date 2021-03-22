package evtq

import (
	"github.com/chuckha/evtq/core"
	"github.com/chuckha/evtq/usecases"
)

type Broker struct {
	Adapter
	UseCases
}

func NewBroker(adapter Adapter, useCases UseCases) *Broker {
	return &Broker{
		Adapter:  adapter,
		UseCases: useCases,
	}
}

func (b *Broker) AddConnectorFromInfo(info *core.ConnectorBuilderInfo) error {
	connector, err := b.Adapter.AddConnectorFromInfo(info)
	if err != nil {
		return err
	}
	return b.UseCases.AddConnector(connector)
}

type Adapter interface {
	AddConnectorFromInfo(info *core.ConnectorBuilderInfo) (*core.Connector, error)
}

type UseCases interface {
	AddConnector(connector *core.Connector) error
}

type connectorBuilder interface {
	BuildConnector(info *core.ConnectorBuilderInfo) (*core.Connector, error)
}

type BrokerAdapter struct {
	connectorBuilder
}

func (b *BrokerAdapter) AddConnectorFromInfo(info *core.ConnectorBuilderInfo) (*core.Connector, error) {
	return b.BuildConnector(info)
}

type BrokerUseCases struct {
	*usecases.ConsumerConnector
}
