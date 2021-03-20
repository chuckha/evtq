package domain

import (
	"testing"
)

func TestDefaultCompressor_Compress(t *testing.T) {
	c := NewDefaultCompressor()
	evt, err := NewEvent("myevt", []byte("my Data"))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.Compress(evt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultCompressor_Compress_TracksMaxOffset(t *testing.T) {
	c := NewDefaultCompressor()
	evt, err := NewEvent("myevt", []byte("my Data"))
	if err != nil {
		t.Fatal(err)
	}
	data, offset, err := c.Compress(evt)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("ensure first offset is same length as data", func(tt *testing.T) {
		if offset != len(data) {
			t.Fatalf("offset (%d) and data (%d) were not the same length", offset, len(data))
		}
	})
	data, offset2, err := c.Compress(evt)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("ensure second offset is twice the size of the first offset", func(tt *testing.T) {
		if offset2 != 2*offset {
			t.Fatalf("second offset (%d) should be twice the first (%d) but it is not", offset2, offset)
		}
	})
}

func TestDefaultDecompressor_Decompress(t *testing.T) {
	d := &DefaultDecompressor{}
	_, err := d.Decompress([]byte("{}"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultOffsetManager(t *testing.T) {
	om := NewDefaultOffsetManager()
	om.Set("one-one", 1)
	om.Set("two-one", 2)
	if om.Get(ProcessorID("one", "one")) != 1 {
		t.Fatal("got the wrong data")
	}
	if om.Get(ProcessorID("two", "one")) != 2 {
		t.Fatal("got the wrong data")
	}
}

type rcv struct {
	name    string
	evtType string
	count   int
}

func newRcv(name, t string) *rcv {
	return &rcv{name, t, 0}
}

func (r *rcv) Receive(_ []byte) error {
	r.count++
	return nil
}

func (r *rcv) Name() string {
	return r.name
}

func (r *rcv) EventTypes() []string {
	return []string{r.evtType}
}

func TestDefaultReceiverManager_AddReceiver(t *testing.T) {
	rm := NewDefaultReceiverManager()
	evtType := "one"
	rcv1 := newRcv("one", evtType)
	if err := rm.AddReceiver(rcv1); err != nil {
		t.Fatal(err)
	}
	rcvs := rm.GetReceiver(evtType)
	if len(rcvs) != 1 {
		t.Fatal("failed to get receivers by event type")
	}
	rm.RemoveReceiver(rcv1)
	rcvs = rm.GetReceiver(evtType)
	if len(rcvs) != 0 {
		t.Fatal("faield to remove receiver properly")
	}
}

func TestDefaultReceiverManager_GetReceiver(t *testing.T) {
	rm := NewDefaultReceiverManager()
	rcvs := rm.GetReceiver("unknown event")
	if len(rcvs) != 0 {
		t.Fatal("got nothing from something")
	}
}

func TestNewDefaultReceiverManager(t *testing.T) {
	rm := NewDefaultReceiverManager()
	rcv1 := newRcv("one", "hello")
	rcv2 := newRcv("one", "other")
	t.Run("cannot add two receivers of the same name", func(tt *testing.T) {
		if err := rm.AddReceiver(rcv1); err != nil {
			tt.Fatal("should be able to add this receiver")
		}
		if err := rm.AddReceiver(rcv2); err == nil {
			tt.Fatal("should not be able to add this receiver")
		}
	})
	rcv3 := newRcv("three", "hello")
	t.Run("should be able to add another receiver of a different name", func(tt *testing.T) {
		if err := rm.AddReceiver(rcv3); err != nil {
			tt.Fatal("should be allowed to add a receiver with a new name")
		}
	})
}
