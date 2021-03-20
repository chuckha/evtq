package graveyard

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type FileFactory interface {
	GetOrCreateFile(evtType string) (io.WriteCloser, error)
}

type DefaultFileFactory struct {
	directory string
}

func NewDefaultFileFactory(directory string) (*DefaultFileFactory, error) {
	ff := &DefaultFileFactory{directory: directory}
	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return ff, errors.WithStack(os.Mkdir(directory, 0700))
	}
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return ff, nil
	}
	return nil, errors.Errorf("%q exists but it is a file and must be a directory", directory)
}

func (d *DefaultFileFactory) GetOrCreateFile(evtType string) (io.WriteCloser, error) {
	info, err := os.Stat(d.fileForEventType(evtType))
	if os.IsNotExist(err) {
		return d.createFile(evtType)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if info.IsDir() {
		return nil, errors.Errorf("need to write to file %q but that is already a directory", d.fileForEventType(evtType))
	}
	f, err := os.OpenFile(d.fileForEventType(evtType), os.O_APPEND|os.O_RDWR, 0644)
	return f, errors.WithStack(err)
}

func (d *DefaultFileFactory) createFile(evtType string) (io.WriteCloser, error) {
	f, err := os.OpenFile(d.fileForEventType(evtType), os.O_CREATE|os.O_RDWR, 0644)
	return f, errors.WithStack(err)
}

func (d *DefaultFileFactory) fileForEventType(evtType string) string {
	return filepath.Join(d.directory, evtType)
}

// EventFileWriter writes a directory with one file-per-type
type EventFileWriter struct {
	OpenFiles   map[string]io.WriteCloser
	FileFactory FileFactory
}

func NewEventFileWriter(directory string) (*EventFileWriter, error) {
	ff, err := NewDefaultFileFactory(directory)
	if err != nil {
		return nil, err
	}
	return &EventFileWriter{
		OpenFiles:   make(map[string]io.WriteCloser),
		FileFactory: ff,
	}, nil
}

func (e *EventFileWriter) WriteEvent(data []byte, evtType string) error {
	var err error
	f, ok := e.OpenFiles[evtType]
	if !ok {
		f, err = e.FileFactory.GetOrCreateFile(evtType)
		if err != nil {
			return err
		}
		e.OpenFiles[evtType] = f
	}
	_, err = f.Write(data)
	return errors.WithStack(err)
}

func (e *EventFileWriter) Shutdown() {
	for name, wc := range e.OpenFiles {
		if err := wc.Close(); err != nil {
			fmt.Println("trouble shutting down the event file writer", name)
		}
	}
}
