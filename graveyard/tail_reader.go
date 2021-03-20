package graveyard

import (
	"io"
	"time"
)

type tailReader struct {
	io.Reader
	closec chan struct{}
}

func newTailReader(r io.Reader, closec chan struct{}) *tailReader {
	return &tailReader{
		Reader: r,
		closec: closec,
	}
}

func (t *tailReader) Read(b []byte) (int, error) {
	ticker := time.NewTicker(10 * time.Millisecond)
	// TODO: i really want this to read then tick...not just start ticking
	for {
		select {
		case <-t.closec:
			ticker.Stop()
			return 0, io.EOF
		case <-ticker.C:
			n, err := t.Reader.Read(b)
			if n > 0 {
				return n, nil
			} else if err != io.EOF {
				return n, err
			}
		}
	}
}
