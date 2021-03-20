package graveyard

import (
	"testing"
)

type rcv struct {
	name  string
	count int
}

func (r *rcv) Receive(_ []byte) error {
	r.count++
	return nil
}

func (r *rcv) Name() string {
	return r.name
}

func (r *rcv) EventTypes() []string {
	return []string{}
}

func TestEventFile(t *testing.T) {
	// f, err := ioutil.TempFile("", "events")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer os.Remove(f.Name())
	//
	// evts := NewEventFile(f)
	// if err := evts.WriteEvent([]byte(`{"one":1,"two":2}`)); err != nil {
	// 	t.Fatal(err)
	// }
	// rcv1 := &rcv{"one", 0}
	// if err := evts.AddReceiver(rcv1, 0); err != nil {
	// 	t.Fatal(err)
	// }
	// if err := evts.WriteEvent([]byte(`{"one":3,"two":5}`)); err != nil {
	// 	t.Fatal(err)
	// }
	// rcv2 := &rcv{"two", 0}
	// if err := evts.AddReceiver(rcv2, 2); err != nil {
	// 	t.Fatal(err)
	// }
	// if err := evts.WriteEvent([]byte(`{"one":4,"two":8}`)); err != nil {
	// 	t.Fatal(err)
	// }
	// if rcv1.count != 3 {
	// 	t.Fatalf("rcv1 should have seen 3 events but was: %d", rcv1.count)
	// }
	// if rcv2.count != 1 {
	// 	t.Fatalf("rcv2 should have seen 1 events but was: %d", rcv2.count)
	// }
	// time.Sleep(50 * time.Millisecond)
	// evts.Close()
}
