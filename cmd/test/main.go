package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

// one buf == one event

type SingleWriterManyReaders interface {
	io.WriteCloser
	NewReader() (io.ReadCloser, error)
}

type swmr struct {
	*os.File
}

type tailReader struct {
	io.ReadCloser
}

func (t *tailReader) Read(b []byte) (int, error) {
	for {
		n, err := t.ReadCloser.Read(b)
		if n > 0 {
			return n, nil
		} else if err != io.EOF {
			return n, err
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (s *swmr) NewReader() (*tailReader, error) {
	f, err := os.Open(s.Name())
	if err != nil {
		return nil, err
	}
	return &tailReader{f}, nil
}

func main() {
	f, err := os.OpenFile("my-test-file", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	mb := &swmr{f}
	// write every 5 seconds
	writeTicker := time.Tick(5 * time.Second)
	go func() {
		for {
			select {
			case <-writeTicker:
				if _, err := mb.Write([]byte("5 seconds has elapsed.\n")); err != nil {
					panic(err)
				}
			}
		}
	}()

	// have three readers reading concurrently
	for i := 0; i < 3; i++ {
		go func(i int) {
			r, err := mb.NewReader()
			if err != nil {
				panic(err)
			}
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				fmt.Printf("[%d]: %s\n", i, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading:", err)
			}
		}(i)
	}

	time.Sleep(5 * time.Second)
	go func(i int) {
		r, err := mb.NewReader()
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			fmt.Printf("[%d]: %s\n", i, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading:", err)
		}
	}(4)

	select {}
}
