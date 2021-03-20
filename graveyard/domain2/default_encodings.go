package domain2

import (
	"encoding/gob"
	"encoding/json"
	"io"
)

type JSONEncDec struct{}
type GobEncDec struct{}

func (j *JSONEncDec) NewEncoder(writer io.Writer) Encoder {
	return json.NewEncoder(writer)
}

func (j *JSONEncDec) NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}

func (g *GobEncDec) NewEncoder(writer io.Writer) Encoder {
	return gob.NewEncoder(writer)
}

func (g *GobEncDec) NewDecoder(reader io.Reader) Decoder {
	return gob.NewDecoder(reader)
}
