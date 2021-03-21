package infrastructure

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/chuckha/evtq/core"
)

type testlog struct{}

func (t *testlog) debug(_ string) {}

type myconnector struct {
	name        string
	failSendErr error
	offsets     []*core.Offset
}

func (m *myconnector) GetName() string {
	return m.name
}

func (m *myconnector) GetOffsets() []*core.Offset {
	return m.offsets
}

func (m *myconnector) SendEvents(e ...*core.Event) error {
	return m.failSendErr
}
func newMyConnector(name string, failSendErr error, offsets []*core.Offset) *myconnector {
	return &myconnector{
		name:        name,
		failSendErr: failSendErr,
		offsets:     offsets,
	}
}

func TestNewConnectorsRegistry(t *testing.T) {
	cr := NewConnectorsRegistry(WithLog(&testlog{}))
	c1 := newMyConnector("c1", errors.New("ffff"), []*core.Offset{{EventType: "test1"}})
	cr.RegisterConnector(c1)
	evt1, err := core.NewEvent("test1", []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	cr.DistributeEvent(evt1)
	if _, ok := cr.connectors[c1.name]; ok {
		t.Errorf("c1 should have been removed from connectors as it failed to send")
	}
	if _, ok := cr.connectorsByEvt["test1"][c1.name]; ok {
		t.Errorf("c1 should have been removed from connectors by event as it failed to send")
	}

	// i f
}
