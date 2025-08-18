package ch04

import (
	"errors"
	"log"
	"net"
	"testing"
	"time"
)

var (
	err        error
	n          int
	maxRetries = 7 // maximum number of retries
)

func TestErrorChecking(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for ; maxRetries > 0; maxRetries-- {
		n, err = conn.Write([]byte("hello world"))
		if err != nil {
			var netError net.Error
			if errors.As(err, &netError) {
				log.Println("temporary error: ", netError)
				time.Sleep(10 * time.Second)
				continue
			}
			t.Fatal(err)
		}
		break
	}

	if maxRetries == 0 {
		t.Error("temporary write failure threshold exceeded")
	}

	log.Printf("wrote %d bytes to %s\n", n, conn.RemoteAddr())
}
