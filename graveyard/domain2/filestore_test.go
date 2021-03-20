package domain2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestFileEventStore_EventsFromOffset(t *testing.T) {
	encdec := &GobEncDec{}
	dir, err := ioutil.TempDir("", "tests")
	if err != nil {
		t.Fatal("failed to create temp dir")
	}
	defer os.RemoveAll(dir)
	store := NewStore(NewFileStore(dir, encdec))
	etype := "abc"
	for i := 0; i < 5; i++ {
		e, err := NewEvent(etype, []byte(fmt.Sprintf("hello world %d", i)))
		if err != nil {
			t.Fatal(fmt.Sprintf("%+v", err))
		}
		if err := store.Write(e); err != nil {
			t.Fatal(fmt.Sprintf("%+v", err))
		}
	}

	evts, err := store.EventsFromOffset(etype, 2)
	if err != nil {
		t.Fatal(fmt.Sprintf("%+v", err))
	}
	//sender
	var buf bytes.Buffer
	sender := &Sender{EncDec: encdec}
	if err := sender.SendEvents(&buf, evts...); err != nil {
		t.Fatal(err)
	}

	handler := &EventCountHandler{}
	go func() {
		if err := NewConsumer(encdec, handler).Consume(&buf); err != nil {
			fmt.Printf("%+v", err)
		}
	}()
	time.Sleep(200 * time.Millisecond)
	if handler.Count != 3 {
		t.Fatalf("did not read all the events, %v", handler.Count)
	}
}

// func TestFileEventStore_EventListener(t *testing.T) {
// 	encdec := &GobEncDec{}
// 	dir, err := ioutil.TempDir("", "tests")
// 	if err != nil {
// 		t.Fatal("failed to create temp dir")
// 	}
// 	defer os.RemoveAll(dir)
// 	filestore := NewFileStore(dir)
// 	store := NewStore(&GobEncoding{}, filestore)
// 	etype := "abc"
// 	reader := store.EventListener(etype)
// 	numEvents := 5
// 	for i := 0; i < numEvents; i++ {
// 		e, err := NewEvent(etype, []byte(fmt.Sprintf("hello world %d", i)))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if err := store.Write(e); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
//
// 	handler := &EventCountHandler{}
// 	go NewConsumer(&GobDecoding{}, handler).Consume(reader)
// 	time.Sleep(120 * time.Millisecond)
// 	if handler.Count != numEvents {
// 		t.Fatalf("did not read all the events, %v", handler.Count)
// 	}
// }
