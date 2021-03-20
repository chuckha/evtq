package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileStore_GetEventsFrom(t *testing.T) {
	out, err := ioutil.TempDir("", "testing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(out)
	fs := NewFileStore(out, &JSONEncDec{})
	evtType := "abcd"
	t.Run("no file should give us no events", func(tt *testing.T) {
		events, err := fs.GetEventsFrom(evtType, 4)
		if err != nil {
			tt.Fatal(err)
		}
		if len(events) != 0 {
			tt.Fatal("expected no output")
		}
	})

	t.Run("a file with no contents should give us no events", func(tt *testing.T) {
		f, err := os.Create(filepath.Join(out, evtType))
		if err != nil {
			tt.Fatal(err)
		}
		if err := f.Close(); err != nil {
			tt.Fatal(err)
		}
		events, err := fs.GetEventsFrom(evtType, 4)
		if err != nil {
			t.Fatal(err)
		}
		if len(events) != 0 {
			t.Fatal("expected no output")
		}
	})

	t.Run("a file with some events should give us some events", func(tt *testing.T) {
		evt, err := NewEvent(evtType, []byte("hello world!"))
		if err != nil {
			tt.Fatal(err)
		}
		if err := fs.WriteEvent(evt); err != nil {
			tt.Fatalf("%+v", err)
		}
		events, err := fs.GetEventsFrom(evtType, 0)
		if err != nil {
			tt.Fatalf("%+v", err)
		}
		if len(events) != 1 {
			tt.Fatalf("expected 1 event but got %d", len(events))
		}
	})

	t.Run("a file with a lot of events should allow us to access any subset of events", func(tt *testing.T) {
		for i := 0; i < 9; i++ {
			evt, err := NewEvent(evtType, []byte("hello world!"))
			if err != nil {
				tt.Fatalf("%+v", err)
			}
			if err := fs.WriteEvent(evt); err != nil {
				tt.Fatalf("%+v", err)
			}
		}
		events, err := fs.GetEventsFrom(evtType, 7)
		if err != nil {
			tt.Fatalf("%+v", err)
		}
		if len(events) != 3 {
			tt.Fatalf("expected 3 event but got %d", len(events))
		}
	})

	t.Run("a lot of events should be able to be read from some point", func(tt *testing.T) {
		for i := 0; i < 890; i++ {
			evt, err := NewEvent(evtType, []byte("hello world!"))
			if err != nil {
				tt.Fatalf("%+v", err)
			}
			if err := fs.WriteEvent(evt); err != nil {
				tt.Fatalf("%+v", err)
			}
		}
		events, err := fs.GetEventsFrom(evtType, 100)
		if err != nil {
			tt.Fatalf("%+v", err)
		}
		if len(events) != 800 {
			tt.Fatalf("expected 800 event but got %d", len(events))
		}
	})
}

func BenchmarkFileStore_WriteEvent(b *testing.B) {
	out, err := ioutil.TempDir("", "testing")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(out)
	fs := NewFileStore(out, &GobEncDec{})
	evtType := "abcd"
	evt, err := NewEvent(evtType, []byte(`{"onefield": "hello", "twoField": 33345'`))
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		if err := fs.WriteEvent(evt); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFileStore_GetEventsFrom(b *testing.B) {
	out, err := ioutil.TempDir("", "testing")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(out)
	fs := NewFileStore(out, &GobEncDec{})
	evtType := "abcd"
	evt, err := NewEvent(evtType, []byte(`{"onefield": "hello", "twoField": 33345'`))
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
