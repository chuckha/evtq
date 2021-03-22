package infrastructure

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/chuckha/evtq/core"
)

func TestLocalConnector(t *testing.T) {
	lc, err := NewLocalConnector(&core.ConnectorBuilderInfo{
		Name:         "my connector",
		EventTypes:   []string{},
		EncodingType: GOBEncoding,
		Info:         LocalConnectorInfo{},
	}, []*core.Offset{core.NewOffset("none", 1)})
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	readyc := make(chan struct{})
	closec := make(chan struct{})
	donec := make(chan struct{})
	go func(reader io.Reader, readyc, closec, donec chan struct{}) {
		for {
			select {
			case <-closec:
				return
			case <-readyc:
				scanner := bufio.NewScanner(reader)
				if !scanner.Scan() {
					fmt.Println("failed scan")
					return
				}
				o := scanner.Bytes()
				if len(o) > 0 {
					count++
				}
				donec <- struct{}{}
			}
		}
	}(lc.GetReadWriter(), readyc, closec, donec)
	evt, err := core.NewEvent("test", []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	if err := lc.SendEvents(evt); err != nil {
		t.Fatal(err)
	}
	readyc <- struct{}{}
	<-donec
	closec <- struct{}{}
	if count != 1 {
		t.Fatal("expected 1 event to have been consumed")
	}
}

func TestNewTCPConnector(t *testing.T) {
	receivedMessages := 0
	closec := make(chan struct{})
	readyc := make(chan struct{})
	go func(readyc, closec chan struct{}) {
		l, err := net.Listen("tcp", ":4000")
		if err != nil {
			t.Fatal(err)
		}
		for {
			select {
			case <-closec:
				if err := l.Close(); err != nil {
					fmt.Println("an error closing connection", err)
				}
			case readyc <- struct{}{}:
				conn, err := l.Accept()
				if err != nil {
					fmt.Println("an error in accept", err)
					continue
				}
				receivedMessages++
				if err := conn.Close(); err != nil {
					fmt.Println("an error closing connection", err)
				}
			}
		}
	}(readyc, closec)
	<-readyc
	connector, err := NewTCPConnector(&core.ConnectorBuilderInfo{
		Name:         "testing",
		EventTypes:   []string{},
		EncodingType: GOBEncoding,
		Info: TCPConnectorInfo{
			Network: "tcp",
			Address: ":4000",
		},
	}, []*core.Offset{core.NewOffset("test1", 0)})
	if err != nil {
		t.Fatal(err)
	}
	evt, err := core.NewEvent("test", []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	if err := connector.SendEvents(evt); err != nil {
		t.Fatal(err)
	}
	closec <- struct{}{}
	if receivedMessages != 1 {
		t.Fatal("should have received an event but did not")
	}
}
