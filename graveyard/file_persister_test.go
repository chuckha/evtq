package graveyard

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenOrCreateWriter(t *testing.T) {
	evtType := "test-event-type"
	dir, err := ioutil.TempDir("", "testing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	f, err := os.OpenFile(filepath.Join(dir, evtType), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	ew := newEventWriter(f)
	numObjs := 5
	for i := 0; i < numObjs; i++ {
		if err := ew.WriteEvent([]byte("{}")); err != nil {
			t.Fatal(err)
		}
	}
	if err := ew.Close(); err != nil {
		t.Fatal(err)
	}

	p := NewFilePersister(dir)
	if err := p.PersistEvent([]byte("hello"), "my-type"); err != nil {
		t.Fatal(err)
	}
}
