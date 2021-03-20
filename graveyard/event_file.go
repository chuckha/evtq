package graveyard

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/chuckha/evtq/graveyard/domain"
)

// TODO: figure out the objects here and abstract so that index could be a db entry or byte position.
// this is too heavily bound to files and line-by-line reading

type eventReaderFactory struct {
	file       *os.File
	closeChans []chan struct{}
}

func NewEventReaderFactory(file *os.File) *eventReaderFactory {
	return &eventReaderFactory{
		file:       file,
		closeChans: make([]chan struct{}, 0),
	}
}

func (e *eventReaderFactory) Close() {
	for _, c := range e.closeChans {
		c <- struct{}{}
	}
}

func (e *eventReaderFactory) newReader() (io.Reader, error) {
	f, err := os.Open(e.file.Name())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	closec := make(chan struct{}, 1)
	e.closeChans = append(e.closeChans, closec)
	tail := newTailReader(f, closec)
	return tail, nil
}

//
// func AddReceiver(r domain.Receiver, offset int) error {
// 	reader, err := e.newReader()
// 	if err != nil {
// 		return err
// 	}
// 	go receiverRead(r, reader, offset)
// 	return nil
// }

// receiverRead will read line by line (event by event) until it's at the point the receiver wants the data.
// Then it will send line by line to the receiver
func receiverRead(r domain.Receiver, reader io.Reader, offset int) {
	scanner := bufio.NewScanner(reader)
	i := 0
	for scanner.Scan() {
		data := scanner.Bytes()
		if i < offset {
			i++
			continue
		}
		if err := r.Receive(data); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[%s] errored receiving bytes: %v", r.Name(), err)
			return
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[%s] errored after scan finished: %v", r.Name(), err)
		return
	}
}
