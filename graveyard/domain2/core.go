package domain2

type ConsumerBuilderInfo struct {
	Name       string
	EventTypes []string
	Info       interface{}
}

// The writer builder should wrap the writer in an encoder
type ConsumerBuilder interface {
	Build(info *ConsumerBuilderInfo) (*Consumer, error)
}

// type Consumer struct {
// 	io.Writer
//
// 	Name       string
// 	EventTypes []string
// }

// when i receive a new writerbuilder info
// i build a writer and add it to my list of writers
// i must also be reading whatever the writer is writing to, but that's a client concern

// persister that is an io.Writer
// a receiver that is an io.Writer
// a reader that is an io.Reader
