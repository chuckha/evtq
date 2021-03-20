package domain

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestNewEvent(t *testing.T) {
	t.Run("cannot create event with no event type", func(tt *testing.T) {
		_, err := NewEvent("", []byte(""))
		if err == nil {
			t.Fatal("this event is not allowed")
		}
	})
}

func TestNewBroker(t *testing.T) {
	c := NewDefaultCompressor()
	o := &DefaultOffsetManager{}
	r := &DefaultReceiverManager{}
	testcases := []struct {
		name string
		c    Compressor
		o    OffsetManager
		r    ReceiverManager
	}{
		{
			name: "missing compressor",
			c:    nil, o: o, r: r,
		},
		{
			name: "missing offset manager",
			c:    c, o: nil, r: r,
		},
		{
			name: "missing receiver manager",
			c:    c, o: o, r: nil,
		},
		{
			name: "missing all",
			c:    nil, o: nil, r: nil,
		},
	}
	// run success case
	t.Run("can create a broker with all valid inputs", func(tt *testing.T) {
		_, err := NewBroker(c, o, r)
		if err != nil {
			tt.Fatal("broker should have succeeded")
		}
	})

	for _, tc := range testcases {
		t.Run(tc.name, func(tt *testing.T) {
			_, err := NewBroker(tc.c, tc.o, tc.r)
			if err == nil {
				tt.Fatal("broker should have failed")
			}
		})
	}
}

func TestBrokerAccept(t *testing.T) {
	c := NewDefaultCompressor()
	o := NewDefaultOffsetManager()
	r := NewDefaultReceiverManager()
	b, err := NewBroker(c, o, r)
	if err != nil {
		t.Fatal("should have succeeded")
	}
	evt1Type := "my-event"
	evt2Type := "your-event"
	evt1 := mustNewEvent(t, evt1Type, []byte("mine"))
	evt2 := mustNewEvent(t, evt2Type, []byte("yours"))

	t.Run("send to no receivers", func(tt *testing.T) {
		if err := b.Accept(evt1); err != nil {
			tt.Fatal(err)
		}
	})

	rcv1 := newRcv("one", evt1Type)
	if err := b.AddReceiver(rcv1); err != nil {
		t.Fatal(err)
	}

	t.Run("sends an event, but no receivers will hear", func(tt *testing.T) {
		if err := b.Accept(evt2); err != nil {
			tt.Fatal(err)
		}
		if rcv1.count != 0 {
			tt.Fatal("rcv1 got a message when it should not")
		}
	})

	rcv2 := newRcv("two", evt2Type)
	if err := b.AddReceiver(rcv2); err != nil {
		t.Fatal(err)
	}

	t.Run("send an event that one out of two receivers will receive", func(tt *testing.T) {
		if err := b.Accept(evt1); err != nil {
			tt.Fatal("failed to accept event")
		}
		if rcv2.count != 0 {
			tt.Fatal("rcv2 should not have received any data")
		}
		if rcv1.count != 1 {
			tt.Fatalf("rcv1 should have received one event, but got %d", rcv1.count)
		}
	})
	t.Run("test the offset was calculated correctly", func(tt *testing.T) {
		offset := o.Get(ProcessorID(rcv1.name, rcv1.EventTypes()[0]))
		if offset == 0 {
			tt.Fatal("offset was not set appropriately")
		}
	})

	t.Run("manually set a very high offset for rcv1", func(tt *testing.T) {
		o.Set(ProcessorID(rcv1.name, rcv1.EventTypes()[0]), 100)
	})

	t.Run("test the broker skips the evt if the offset is too high", func(tt *testing.T) {
		if err := b.Accept(evt1); err != nil {
			tt.Fatal(err)
		}
		if rcv1.count != 1 {
			tt.Fatal("rcv1 received an event it should not have")
		}
	})
}

type failureReceiver struct {
	name    string
	evtType string
}

func newFailureReceiver(n, e string) *failureReceiver {
	return &failureReceiver{n, e}
}
func (f *failureReceiver) Receive(evt []byte) error {
	return errors.New("an error")
}

func (f *failureReceiver) Name() string {
	return f.name
}

func (f *failureReceiver) EventTypes() []string {
	return []string{f.evtType}
}

func TestBroker_Accept(t *testing.T) {
	c := NewDefaultCompressor()
	o := NewDefaultOffsetManager()
	r := NewDefaultReceiverManager()
	b, err := NewBroker(c, o, r)
	if err != nil {
		t.Fatal("should have succeeded")
	}
	evt1Type := "my-event"
	evt1 := mustNewEvent(t, evt1Type, []byte("mine"))
	rcv1 := newFailureReceiver("whatever", evt1Type)
	if err := b.AddReceiver(rcv1); err != nil {
		t.Fatal(err)
	}
	t.Run("test a failed receive does not update the offset", func(tt *testing.T) {
		if err := b.Accept(evt1); err == nil {
			t.Fatal("expected an error but didn't get one")
		}
		if b.offsetManager.Get(ProcessorID(rcv1.Name(), rcv1.EventTypes()[0])) != 0 {
			t.Fatal("offset updated but it should not have")
		}
	})
}

func mustNewEvent(t *testing.T, evtType string, data []byte) *Event {
	evt1, err := NewEvent(evtType, data)
	if err != nil {
		t.Fatal(err)
	}
	return evt1
}

type brokenCompressor struct{}

func (b brokenCompressor) Compress(evt *Event) ([]byte, int, error) {
	return nil, 0, errors.New("an error")
}

func TestBroker_Accept_Bad_Compressor(t *testing.T) {
	c := &brokenCompressor{}
	o := NewDefaultOffsetManager()
	r := NewDefaultReceiverManager()
	b, err := NewBroker(c, o, r)
	if err != nil {
		t.Fatal("should have succeeded")
	}
	e1 := mustNewEvent(t, "banana", []byte("hi"))
	if err := b.Accept(e1); err == nil {
		t.Fatal("should have raised an error but did not")
	}
}

func BenchmarkBroker_Accept(b *testing.B) {
	c := NewDefaultCompressor()
	o := NewDefaultOffsetManager()
	r := NewDefaultReceiverManager()
	broker, err := NewBroker(c, o, r)
	if err != nil {
		b.Fatal("should have succeeded")
	}
	for i := 0; i < 1000; i++ {
		rcv := newRcv(fmt.Sprintf("rcv%d", i), fmt.Sprintf("evt%d", i%10))
		if err := broker.AddReceiver(rcv); err != nil {
			b.Fatal("should not fail adding receiver")
		}
	}
	evts := make([]*Event, 10)
	for i := 0; i < 10; i++ {
		evts[i], err = NewEvent(fmt.Sprintf("evt%d", i), bytes.Repeat([]byte("b"), 100))
	}
	for i := 0; i < b.N; i++ {
		if err := broker.Accept(evts[i%10]); err != nil {
			b.Error(err)
		}
	}
}
