package graveyard

//
// import "errors"
//
// // consumer wraps up processors in a convenient to use fashion.
// // the consumer gets an Event and can route the Event to the associated processor.
// type consumer struct {
// 	// types maps an Event type this consumer is interested in to the last known offset received.
// 	processors map[string]Processor
// }
//
// // receiver handles receiving messages from the broker and dispatching them to their consumer.
// type receiver struct {
// 	name     string
// 	consumer *consumer
// 	Decompressor
// }
//
// func NewReceiver(name string, c *consumer, d Decompressor) (*receiver, error) {
// 	if name == "" {
// 		return nil, errors.New("receiver name cannot be empty")
// 	}
// 	if d == nil {
// 		return nil, errors.New("decompressor must exist")
// 	}
// 	return &receiver{
// 		name:         name,
// 		consumer:     c,
// 		Decompressor: d,
// 	}, nil
// }
//
// func (e *receiver) receiveEvent(data []byte, offset int) error {
// 	evt, err := e.Decompress(data)
// 	if err != nil {
// 		return err
// 	}
// 	return e.consumer.process(evt, offset)
// }
//
// func (r *receiver) EventTypes() []eventType {
// 	return r.consumer.eventTypes()
// }
//
// func NewConsumer(processors ...Processor) (*consumer, error) {
// 	c := &consumer{
// 		processors: make(map[string]Processor),
// 	}
// 	for _, processor := range processors {
// 		if err := c.registerProcessor(processor); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return c, nil
// }
//
// func (c *consumer) registerProcessor(p Processor) error {
// 	if _, ok := c.processors[p.EventType()]; ok {
// 		return errors.New("cannot overwrite an existing processor, create a new consumer")
// 	}
// 	c.processors[p.EventType()] = p
// 	return nil
// }
//
// func (c *consumer) process(evt *Event, offset int) error {
// 	processor := c.processors[evt.Type]
// 	if processor.Offset() >= offset {
// 		return nil
// 	}
// 	if err := processor.Process(evt); err != nil {
// 		return err
// 	}
// 	processor.SetOffset(offset)
// 	return nil
// }
//
// func (c *consumer) eventTypes() []eventType {
// 	out := []eventType{}
// 	for key := range c.processors {
// 		out = append(out, eventType(key))
// 	}
// 	return out
// }
