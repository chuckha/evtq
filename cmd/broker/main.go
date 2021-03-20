package main

import (
	"github.com/chuckha/evtq/graveyard"
	"github.com/chuckha/evtq/graveyard/usecases"
)

func main() {
	broker := &graveyard.Broker{
		EventAdder: &usecases.EventAdder{
			Compressor:     nil,
			EventPersister: nil,
			ReceiverFinder: nil,
			OffsetUpdater:  nil,
		},
		ReceiverAdder:   nil,
		ReceiverRemover: nil,
	}
}
