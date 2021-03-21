package infrastructure

import (
	"testing"

	"github.com/chuckha/evtq/core"
)

func TestMemoryStore_GetEventsFrom(t *testing.T) {
	ms := NewMemoryStore()
	evtType := "abcd"
	t.Run("an empty memory store should return nothing", func(tt *testing.T) {
		events, err := ms.GetEventsFrom(evtType, 4)
		if err != nil {
			tt.Fatal(err)
		}
		if len(events) != 0 {
			tt.Fatal("expected no output")
		}
	})

	t.Run("some events should give us some events", func(tt *testing.T) {
		evt, err := core.NewEvent(evtType, []byte("hello world!"))
		if err != nil {
			tt.Fatal(err)
		}
		if err := ms.WriteEvent(evt); err != nil {
			tt.Fatalf("%+v", err)
		}
		events, err := ms.GetEventsFrom(evtType, 0)
		if err != nil {
			tt.Fatalf("%+v", err)
		}
		if len(events) != 1 {
			tt.Fatalf("expected 1 event but got %d", len(events))
		}
	})

	t.Run("a file with a lot of events should allow us to access any subset of events", func(tt *testing.T) {
		for i := 0; i < 9; i++ {
			evt, err := core.NewEvent(evtType, []byte("hello world!"))
			if err != nil {
				tt.Fatalf("%+v", err)
			}
			if err := ms.WriteEvent(evt); err != nil {
				tt.Fatalf("%+v", err)
			}
		}
		events, err := ms.GetEventsFrom(evtType, 7)
		if err != nil {
			tt.Fatalf("%+v", err)
		}
		if len(events) != 3 {
			tt.Fatalf("expected 3 event but got %d", len(events))
		}
	})

	t.Run("a lot of events should be able to be read from some point", func(tt *testing.T) {
		for i := 0; i < 890; i++ {
			evt, err := core.NewEvent(evtType, []byte("hello world!"))
			if err != nil {
				tt.Fatalf("%+v", err)
			}
			if err := ms.WriteEvent(evt); err != nil {
				tt.Fatalf("%+v", err)
			}
		}
		events, err := ms.GetEventsFrom(evtType, 100)
		if err != nil {
			tt.Fatalf("%+v", err)
		}
		if len(events) != 800 {
			tt.Fatalf("expected 800 event but got %d", len(events))
		}
	})
}

func BenchmarkMemoryStore_WriteEvent(b *testing.B) {
	fs := NewMemoryStore()
	evtType := "abcd"
	evt, err := core.NewEvent(evtType, []byte(`{"onefield": "hello", "twoField": 33345'`))
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		if err := fs.WriteEvent(evt); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMemoryStore_GetEventsFrom(b *testing.B) {
	fs := NewMemoryStore()
	evtType := "abcd"
	evt, err := core.NewEvent(evtType, []byte(`{"onefield": "hello", "twoField": 33345'`))
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 10000; i++ {
		if err := fs.WriteEvent(evt); err != nil {
			b.Fatal(err)
		}
	}
	for i := 0; i < b.N; i++ {
		_, err := fs.GetEventsFrom(evtType, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}
