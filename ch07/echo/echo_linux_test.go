//go:build linux

package echo

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestEchoServerUnixPacket(t *testing.T) {
	dir, err := os.MkdirTemp("", "echo_unixpacket")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			t.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
	rAddr, err := streamingEchoServer(ctx, "unixpacket", socket)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	err = os.Chmod(socket, os.ModeSocket|066)
	if err != nil {
		t.Fatal(err)
	}

	// unixpacket is session oriented, so use net.Dial
	// like a unix socket
	conn, err := net.Dial("unixpacket", rAddr.String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	// Write 3 "ping" messages
	msg := []byte("ping")
	for range 3 {
		_, err = conn.Write(msg)
		if err != nil {
			t.Fatal(err)
		}
	}

	// But you read packets individually receiving a packet each time
	// like a datagram server
	buf := make([]byte, 1024)
	for range 3 {
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(msg, buf[:n]) {
			t.Errorf("expected reply %q; actual reply %q", msg, buf[:n])
		}
	}

	// unixpacket discards unrequested data in each datagram
	for range 3 {
		_, err = conn.Write(msg)
		if err != nil {
			t.Fatal(err)
		}
	}

	buf = make([]byte, 2) // only read the first 2 bytes of each reply
	for range 3 {
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(msg[:2], buf[:n]) {
			t.Errorf("expected reply %q; actual reply %q", msg[:2], buf[:n])
		}
	}
}
