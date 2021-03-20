package graveyard

import (
	"github.com/chuckha/evtq/graveyard/usecases"
)

type Broker struct {
	*usecases.EventAdder
	*usecases.ReceiverAdder
	*usecases.ReceiverRemover
}
