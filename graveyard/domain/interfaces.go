package domain

// Decompressor extracts the Event from the input bytes
// TODO: consider an io.Reader
type Decompressor interface {
	Decompress(b []byte) (*Event, error)
}

// Compressor compresses the Event into bytes.
// Decompressor needs to be the inverse of this implementation.
type Compressor interface {
	Compress(evt *Event) ([]byte, int, error)
}

// Processor can process a single message type
type Processor interface {
	Process(e *Event) error
	EventType() string
	Offset() int
	SetOffset(int)
}

// Broker interfaces

// Receiver can receive a []byte representation of an Event
type Receiver interface {
	Receive(evt []byte) error
	Name() string
	EventTypes() []string
}

// OffsetManager manages offsets of various processors
type OffsetManager interface {
	Set(pid PID, offset int)
	Get(pid PID) int
}

// ReceiverManager manages the complex relationships of the receivers
type ReceiverManager interface {
	AddReceiver(rcv Receiver) error
	GetReceiver(evtType string) []Receiver
	RemoveReceiver(rcv Receiver)
}
