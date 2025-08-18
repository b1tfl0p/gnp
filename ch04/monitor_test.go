package ch04

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"testing"
)

// Monitor embeds a log.Logger meant for logging network traffic.
type Monitor struct {
	*log.Logger
}

// Write implements the io.Writer interface
func (m *Monitor) Write(p []byte) (int, error) {
	return len(p), m.Output(2, string(p))
}

func TestMonitor(t *testing.T) {
	monitor := &Monitor{Logger: log.New(os.Stdout, "monitor: ", 0)}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		monitor.Fatal(err)
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()

		b := make([]byte, 1024)
		r := io.TeeReader(conn, monitor)

		n, err := r.Read(b)
		if err != nil && !errors.Is(err, io.EOF) {
			monitor.Println(err)
			t.Error(err)
		}

		w := io.MultiWriter(conn, monitor)

		_, err = w.Write(b[:n]) // echo the message
		if err != nil && !errors.Is(err, io.EOF) {
			monitor.Println(err)
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		monitor.Fatal(err)
		t.Fatal(err)
	}

	_, err = conn.Write([]byte("Test\n"))
	if err != nil {
		monitor.Fatal(err)
		t.Fatal(err)
	}

	_ = conn.Close()
	<-done
}
