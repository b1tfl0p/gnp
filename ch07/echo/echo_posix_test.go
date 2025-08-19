//go:build darwin || linux

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

func TestEchoServerUnixDatagram(t *testing.T) {
	dir, err := os.MkdirTemp("", "echo_unixgram")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			t.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	sSocket := filepath.Join(dir, fmt.Sprintf("s%d.sock", os.Getpid()))
	serverAddr, err := datagramEchoServer(ctx, "unixgram", sSocket)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	err = os.Chmod(sSocket, os.ModeSocket|0622)
	if err != nil {
		t.Fatal(err)
	}

	// Unlike unix socket, use net.ListenPacket for datagram-like behavior
	cSocket := filepath.Join(dir, fmt.Sprintf("c%d.sock", os.Getpid()))
	client, err := net.ListenPacket("unixgram", cSocket)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	err = os.Chmod(cSocket, os.ModeSocket|0622)
	if err != nil {
		t.Fatal(err)
	}

	// Write 3 "ping" messages
	msg := make([]byte, 1024)
	for range 3 {
		_, err = client.WriteTo(msg, serverAddr)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Read 3 packets individually
	buf := make([]byte, 1024)
	for range 3 {
		n, addr, err := client.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}

		if addr.String() != serverAddr.String() {
			t.Fatalf("received reply from %q instead of %q", addr, serverAddr)
		}

		if !bytes.Equal(msg, buf[:n]) {
			t.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
		}
	}
}
