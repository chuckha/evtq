package graveyard

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestTailReader_Read(t *testing.T) {
	closec := make(chan struct{})
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "hello world.\n")

	tr := newTailReader(&buf, closec)
	scanner := bufio.NewScanner(tr)
	go func() { time.Sleep(20 * time.Millisecond); closec <- struct{}{} }()
	for scanner.Scan() {
		if scanner.Text() != "hello world." {
			t.Fatal("failed to scan a line")
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal("scanner error", err)
	}
}
