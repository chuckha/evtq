package main

import (
	"github.com/chuckha/evtq"
	"github.com/chuckha/evtq/core"
	"github.com/chuckha/evtq/infrastructure"
	"github.com/chuckha/evtq/usecases"
)

func main() {
	adapter := &evtq.BrokerAdapter{}
	connectors := infrastructure.NewConnectorsRegistry()
	eg := infrastructure.NewMemoryStore()
	or := infrastructure.NewOffsetRepository()
	connAdder := usecases.NewConnectorAdder(connectors, eg, or)
	evtAdder := usecases.NewEventAdder(eg, connectors)
	ucs := &evtq.BrokerUseCases{
		ConnectorAdder: connAdder,
		EventAdder:     evtAdder,
	}
	b := evtq.NewBroker(adapter, ucs)
	// TODO: add your application here and use b as the event bus
	b.AddConnectorFromInfo(&core.ConnectorBuilderInfo{
		Name:         "my-connector",
		EventTypes:   []string{"event-type-1", "event-type-2", "event-type-3"},
		EncodingType: infrastructure.JSONEncoding, // also supports Gob, see infrastructure/encdec.go
		Info:         &infrastructure.LocalConnectorInfo{},
	})
}
