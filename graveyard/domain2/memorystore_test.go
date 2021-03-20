package domain2

// func TestMemoryEventStore_EventsFromOffset(t *testing.T) {
// 	inMem := NewMemoryEventStore()
// 	store := NewStore(&GobEncoding{}, inMem)
// 	etype := "abc"
// 	for i := 0; i < 5; i++ {
// 		e, err := NewEvent(etype, []byte(fmt.Sprintf("hello world %d", i)))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if err := store.Write(e); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
//
// 	reader, err := store.EventsFromOffset(etype, 2)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	handler := &EventCountHandler{}
// 	go func() {
// 		if err := NewConsumer(&GobDecoding{}, handler).Consume(reader); err != nil {
// 			fmt.Printf("%+v", err)
// 		}
// 	}()
// 	time.Sleep(120 * time.Millisecond)
// 	if handler.Count != 3 {
// 		t.Fatalf("did not read all the events, %v", handler.Count)
// 	}
// }
//
// func TestMemoryEventStore_EventListener(t *testing.T) {
// 	inMem := NewMemoryEventStore()
// 	store := NewStore(&GobEncoding{}, inMem)
// 	etype := "abc"
// 	reader := store.EventListener(etype)
// 	numEvents := 5
// 	for i := 0; i < numEvents; i++ {
// 		e, err := NewEvent(etype, []byte(fmt.Sprintf("hello world %d", i)))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if err := store.Write(e); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
//
// 	handler := &EventCountHandler{}
// 	go NewConsumer(&GobDecoding{}, handler).Consume(reader)
// 	time.Sleep(120 * time.Millisecond)
// 	if handler.Count != numEvents {
// 		t.Fatalf("did not read all the events, %v", handler.Count)
// 	}
// }
