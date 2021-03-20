package domain2

import (
	"io"

	"github.com/pkg/errors"
)

type Sender struct {
	EncDec EncDec
}

func (s *Sender) SendEvents(writer io.Writer, events ...*Event) error {
	encoder := s.EncDec.NewEncoder(writer)
	for _, e := range events {
		if err := encoder.Encode(e); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
