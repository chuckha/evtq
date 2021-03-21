package infrastructure

import (
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/chuckha/evtq/core"
)

type JSONEncDec struct{}
type GobEncDec struct{}

func (j *JSONEncDec) NewEncoder(writer io.Writer) core.Encoder {
	return json.NewEncoder(writer)
}

func (j *JSONEncDec) NewDecoder(reader io.Reader) core.Decoder {
	return json.NewDecoder(reader)
}

func (g *GobEncDec) NewEncoder(writer io.Writer) core.Encoder {
	return gob.NewEncoder(writer)
}

func (g *GobEncDec) NewDecoder(reader io.Reader) core.Decoder {
	return gob.NewDecoder(reader)
}

const (
	JSONEncoding core.EncodingType = "json"
	GOBEncoding  core.EncodingType = "gob"
)

func EncDecFactory(encodingType core.EncodingType) core.EncDec {
	switch encodingType {
	case JSONEncoding:
		return &JSONEncDec{}
	case GOBEncoding:
		return &GobEncDec{}
	default:
		return &JSONEncDec{}
	}
}
