package usecases

import (
	"io"
	"testing"

	"github.com/chuckha/evtq/core"
	"github.com/chuckha/evtq/infrastructure"
)

type testlog struct{}

func (t *testlog) Debug(_ string) {}

type tc struct {
	eventsSeen int
	offsets    map[string]*core.Offset
}

func (t *tc) GetName() string {
	return "test connector"
}

func (t *tc) GetOffsets() map[string]*core.Offset {
	return t.offsets
}

func (t *tc) SendEvents(evts ...*core.Event) error {
	for range evts {
		t.eventsSeen++
	}
	return nil
}

func (t *tc) GetReader() io.Reader {
	return nil
}

func TestConsumerConnector_AddConnector(t *testing.T) {
	cr := infrastructure.NewConnectorsRegistry(infrastructure.WithLog(&testlog{}))
	eg := infrastructure.NewMemoryStore()
	or := infrastructure.NewOffsetRepository()
	cc := &ConnectorAdder{
		connectors: cr,
		events:     eg,
		offsets:    or,
	}
	t.Run("ensure events can be distributed after a connector is registered", func(tt *testing.T) {
		mytc := &tc{
			offsets: map[string]*core.Offset{
				"hello": {"hello", 0},
			},
		}
		if _, err := cc.AddConnector(mytc); err != nil {
			tt.Fatal(err)
		}
		n, err := core.NewEvent("hello", []byte("hello world"))
		if err != nil {
			tt.Fatal(err)
		}
		cr.DistributeEvent(n)
		if mytc.eventsSeen != 1 {
			tt.Fatalf("should have seen 1 event but saw %d", mytc.eventsSeen)
		}
	})

	t.Run("ensure a new connector gets all missing events", func(tt *testing.T) {
		mytc := &tc{
			offsets: map[string]*core.Offset{
				"banana": {"banana", 4},
			},
		}
		n, err := core.NewEvent("banana", []byte("hello world"))
		if err != nil {
			tt.Fatal(err)
		}
		for i := 0; i < 10; i++ {
			if err := eg.WriteEvent(n); err != nil {
				tt.Fatal(err)
			}
		}
		if _, err := cc.AddConnector(mytc); err != nil {
			tt.Fatal(err)
		}
		if mytc.eventsSeen != 6 {
			tt.Fatalf("should have seen 6 event but saw %d", mytc.eventsSeen)
		}
	})
}
