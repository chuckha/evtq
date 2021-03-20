package graveyard

import (
	"fmt"
	"testing"

	"github.com/chuckha/evtq/graveyard/domain"
)

func TestEventFileWriter_WriteEvent(t *testing.T) {
	fw, err := NewEventFileWriter("testing")
	if err != nil {
		t.Fatal(err)
	}
	dc := domain.NewDefaultCompressor()
	for i := 0; i < 10; i++ {
		evt, err := domain.NewEvent(fmt.Sprintf("test-%d-evt", i%10), []byte(fmt.Sprintf("hello world %d", i)))
		if err != nil {
			t.Fatal(err)
		}
		data, _, err := dc.Compress(evt)
		if err != nil {
			t.Fatal(err)
		}
		if err := fw.WriteEvent(data, evt.Type); err != nil {
			t.Fatal(err)
		}
	}
	fw.Shutdown()
}
