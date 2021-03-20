package graveyard

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func fileForEventType(directory, evtType string) string {
	return filepath.Join(directory, evtType)
}

type filePersister struct {
	directory  string
	eventFiles map[string]*eventWriter
}

func NewFilePersister(directory string) *filePersister {
	return &filePersister{
		directory:  directory,
		eventFiles: map[string]*eventWriter{},
	}
}

func (f *filePersister) PersistEvent(data []byte, eventType string) error {
	evtf, err := f.getEventFile(eventType)
	if err != nil {
		return err
	}
	return evtf.WriteEvent(data)
}

func (f *filePersister) getEventFile(eventType string) (*eventWriter, error) {
	var err error
	writer, ok := f.eventFiles[eventType]
	if ok {
		return writer, nil
	}
	file, err := f.openOrCreateFile(eventType)
	if err != nil {
		return nil, err
	}
	evtf := newEventWriter(file)
	f.eventFiles[eventType] = evtf
	return evtf, nil
}

func (f *filePersister) openOrCreateFile(evtType string) (*os.File, error) {
	info, err := os.Stat(fileForEventType(f.directory, evtType))
	if os.IsNotExist(err) {
		return f.createFile(evtType)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if info.IsDir() {
		return nil, errors.Errorf("need to write to file %q but that is already a directory", fileForEventType(f.directory, evtType))
	}
	file, err := os.OpenFile(fileForEventType(f.directory, evtType), os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return file, nil
}

func (f *filePersister) createFile(evtType string) (*os.File, error) {
	file, err := os.OpenFile(fileForEventType(f.directory, evtType), os.O_CREATE|os.O_RDWR, 0644)
	return file, errors.WithStack(err)
}

// eventWriter wraps a file to expose a WriteEvent function
type eventWriter struct {
	*os.File
}

func newEventWriter(file *os.File) *eventWriter {
	return &eventWriter{file}
}

func (e *eventWriter) WriteEvent(data []byte) error {
	_, err := e.Write(append(data, []byte("\n")...))
	return errors.WithStack(err)
}
