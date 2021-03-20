package core

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type FileStore struct {
	directory string
	encdec    EncDec
	files     map[string]*os.File
	encoders  map[string]Encoder
}

func NewFileStore(dir string, encdec EncDec) *FileStore {
	return &FileStore{
		directory: dir,
		encdec:    encdec,
		files:     map[string]*os.File{},
		encoders:  map[string]Encoder{},
	}
}

func (f *FileStore) fileName(etype string) string {
	return filepath.Join(f.directory, etype)
}

func (f *FileStore) WriteEvent(e *Event) error {
	etype := e.EventType
	// ensure writer exists
	encoder, ok := f.encoders[etype]
	if !ok {
		file, err := os.OpenFile(f.fileName(etype), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return errors.WithStack(err)
		}
		encoder = f.encdec.NewEncoder(file)
		f.files[etype] = file
		f.encoders[etype] = encoder
	}
	return errors.WithStack(encoder.Encode(e))
}

func (f *FileStore) GetEventsFrom(eventType string, eventNumber int) ([]*Event, error) {
	_, err := os.Stat(f.fileName(eventType))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// get a new file
	file, err := os.Open(f.fileName(eventType))
	if err != nil {
		return nil, err
	}
	dec := f.encdec.NewDecoder(file)
	out := []*Event{}
	i := 0
	for {
		evt := &Event{}
		err := dec.Decode(evt)
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if i < eventNumber {
			i++
			continue
		}
		i++
		out = append(out, evt)
	}
}
