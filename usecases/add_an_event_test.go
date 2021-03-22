package usecases

import (
	"testing"

	"github.com/chuckha/evtq/core"
	"github.com/chuckha/evtq/infrastructure"
)

func TestEventAdder_AddEvent(t *testing.T) {
	ew := infrastructure.NewMemoryStore()
	cr := infrastructure.NewConnectorsRegistry(infrastructure.WithLog(&testlog{}))
	ea := &EventAdder{
		ew: ew,
		ed: cr,
	}

	t.Run("should be able to accept events even with nothing connected", func(tt *testing.T) {
		evt, err := core.NewEvent("zelda", []byte("hey, listen!"))
		if err != nil {
			t.Fatal(err)
		}
		if err := ea.AddEvent(evt); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should dispatch events to connected listeners when new events arrive", func(tt *testing.T) {
		mytc := &tc{
			eventsSeen: 0,
			offsets: []*core.Offset{
				{"link", 0},
			},
		}
		evt, err := core.NewEvent("link", []byte("..."))
		if err != nil {
			t.Fatal(err)
		}
		cr.RegisterConnector(mytc) // manually register
		for i := 0; i < 10; i++ {
			if err := ea.AddEvent(evt); err != nil {
				t.Fatal(err)
			}
		}
		if mytc.eventsSeen != 10 {
			t.Fatal("should have seen 10 but didn't")
		}
	})
}
