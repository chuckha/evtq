package domain2

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type FileStore struct {
	directory string
	counts    map[string]int
	writers   map[string]*os.File
	EncDec    EncDec
}

func NewFileStore(dir string, encdec EncDec) *FileStore {
	return &FileStore{
		directory: dir,
		counts:    map[string]int{},
		writers:   map[string]*os.File{},
		EncDec:    encdec,
	}
}

func (f *FileStore) fileName(etype string) string {
	return filepath.Join(f.directory, etype)
}

func (f *FileStore) getEventsFromOffset(etype string, offset int) ([]*Event, error) {
	file := f.writers[etype]
	reader, err := os.Open(file.Name())
	if err != nil {
		return nil, err
	}
	out := make([]*Event, 0)
	i := 0
	for {
		decoder := f.EncDec.NewDecoder(reader)
		evt := &Event{}
		err := decoder.Decode(evt)
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if i < offset {
			i++
			continue
		}
		out = append(out, evt)
		i++
	}
}

func (f *FileStore) getEventCount(etype string) int {
	return f.counts[etype]
}

func (f *FileStore) write(evt *Event) error {
	etype := evt.Type()
	// ensure writer exists
	var err error
	writer, ok := f.writers[etype]
	if !ok {
		writer, err = os.OpenFile(f.fileName(etype), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return errors.WithStack(err)
		}
		f.writers[etype] = writer
	}
	encoder := f.EncDec.NewEncoder(writer)
	if err := encoder.Encode(evt); err != nil {
		return errors.WithStack(err)
	}
	f.counts[etype]++
	return nil
}
