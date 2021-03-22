package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/chuckha/evtq"
	"github.com/chuckha/evtq/core"
	"github.com/chuckha/evtq/infrastructure"
	"github.com/chuckha/evtq/usecases"
)

func main() {
	eg := infrastructure.NewMemoryStore()
	or := infrastructure.NewOffsetRepository()
	adapter := &evtq.BrokerAdapter{
		ConnectorBuilder: &infrastructure.ConnectorBuilder{
			OffsetRepo: or,
		},
	}
	connectors := infrastructure.NewConnectorsRegistry()
	connAdder := usecases.NewConnectorAdder(connectors, eg, or)
	evtAdder := usecases.NewEventAdder(eg, connectors)
	ucs := &evtq.BrokerUseCases{
		ConnectorAdder: connAdder,
		EventAdder:     evtAdder,
	}
	b := evtq.NewBroker(adapter, ucs)
	// TODO: add your application here and use b as the event bus
	reader, err := b.AddConnector(&core.ConnectorBuilderInfo{
		Name:         "my-connector",
		EventTypes:   []string{"event-type-1", "event-type-2", "event-type-3"},
		EncodingType: infrastructure.JSONEncoding, // also supports Gob, see infrastructure/encdec.go
		Info:         infrastructure.LocalConnectorInfo{},
	})
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	ticker2 := time.Tick(150 * time.Millisecond)
	ticker := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-ticker2:
			if err := b.AddEvent("event-type-1", []byte("hello ticker")); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		case <-ticker:
			decoder := json.NewDecoder(reader)
			evt := &core.Event{}
			err := decoder.Decode(evt)
			if err == io.EOF {
				continue
			}
			if err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			fmt.Println(evt)
		}
	}
}
