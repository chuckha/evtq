package graveyard

import (
	"testing"

	"github.com/chuckha/evtq/graveyard/domain"
)

type counterReceiver struct {
	messagesReceived int
	name             string
}

func (c *counterReceiver) Receive(_ []byte) error {
	c.messagesReceived++
	return nil
}

func (c *counterReceiver) Name() string {
	return c.name
}

func (c *counterReceiver) EventTypes() []string {
	return []string{"evt1", "evt2"}
}

func TestConsumersCanReadFromOffsetZero(t *testing.T) {
	c := domain.NewDefaultCompressor()
	o := domain.NewDefaultOffsetManager()
	r := domain.NewDefaultReceiverManager()
	b, err := domain.NewBroker(c, o, r)
	if err != nil {
		t.Fatal("should have succeeded")
	}
	evt1, err := domain.NewEvent("evt1", []byte("hi world"))
	if err != nil {
		t.Fatal(err)
	}

	// receive a few messages
	for i := 0; i < 4; i++ {
		if err := b.Accept(evt1); err != nil {
			t.Fatal(err)
		}
	}

	rcv2 := &counterReceiver{name: "second"}
	// connect receiver
	if err := b.AddReceiver(rcv2); err != nil {
		t.Fatal(err)
	}

	// ensure they got all the messages
	t.Run("test all messages received", func(tt *testing.T) {
		if rcv2.messagesReceived != 4 {
			tt.Fatalf("expected 4 messages but counted %d", rcv2.messagesReceived)
		}
	})
}
