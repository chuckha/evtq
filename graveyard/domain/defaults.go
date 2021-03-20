package domain

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type MaxOffsetTracker interface {
	// Update takes a new number and returns the old number + new
	Update(evtType string, length int) int
}

type DefaultMaxOffsetTracker struct {
	maxOffsets map[string]int
}

func NewDefaultMaxOffsetTracker() *DefaultMaxOffsetTracker {
	return &DefaultMaxOffsetTracker{maxOffsets: map[string]int{}}
}

func (d *DefaultMaxOffsetTracker) Update(evtType string, length int) int {
	val := length + d.maxOffsets[evtType]
	d.maxOffsets[evtType] = val
	return val
}

type DefaultCompressor struct {
	MaxOffsetTracker
}

func NewDefaultCompressor() *DefaultCompressor {
	return &DefaultCompressor{
		NewDefaultMaxOffsetTracker(),
	}
}

func (c DefaultCompressor) Compress(evt *Event) ([]byte, int, error) {
	o, err := json.Marshal(evt)
	if err != nil {
		return nil, 0, err
	}
	return o, c.Update(evt.Type, len(o)), nil
}

type DefaultDecompressor struct{}

func (d DefaultDecompressor) Decompress(b []byte) (*Event, error) {
	evt := &Event{}
	if err := json.Unmarshal(b, &evt); err != nil {
		return nil, err
	}
	return evt, nil
}

type DefaultOffsetManager struct {
	offsets map[PID]int
}

func NewDefaultOffsetManager() *DefaultOffsetManager {
	return &DefaultOffsetManager{map[PID]int{}}
}

// processor id == rcvName+evtType
func (o *DefaultOffsetManager) Set(pid PID, offset int) {
	o.offsets[pid] = offset
}
func (o *DefaultOffsetManager) Get(pid PID) int {
	return o.offsets[pid]
}

type DefaultReceiverManager struct {
	receivers      map[string]Receiver
	receiversByEvt map[string][]Receiver
}

func NewDefaultReceiverManager() *DefaultReceiverManager {
	return &DefaultReceiverManager{
		receivers:      map[string]Receiver{},
		receiversByEvt: map[string][]Receiver{},
	}
}

func (r *DefaultReceiverManager) AddReceiver(rcv Receiver) error {
	if _, ok := r.receivers[rcv.Name()]; ok {
		return errors.New("cannot overwirte an existing receiver, remove first")
	}
	r.receivers[rcv.Name()] = rcv
	for _, evtType := range rcv.EventTypes() {
		r.receiversByEvt[evtType] = append(r.receiversByEvt[evtType], rcv)
	}
	return nil
}

func (r *DefaultReceiverManager) GetReceiver(evtType string) []Receiver {
	rcvs, ok := r.receiversByEvt[evtType]
	if !ok {
		return []Receiver{}
	}
	return rcvs
}

func (r *DefaultReceiverManager) RemoveReceiver(rcv Receiver) {
	// remove from receivers
	delete(r.receivers, rcv.Name())
	// remove every from the map
	for _, evtt := range rcv.EventTypes() {
		for i, rcver := range r.receiversByEvt[evtt] {
			if rcver.Name() == rcv.Name() {
				r.receiversByEvt[evtt] = append(r.receiversByEvt[evtt][:i], r.receiversByEvt[evtt][i+1:]...)
				break
			}
		}
	}
}
