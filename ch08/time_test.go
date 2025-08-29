package main

import (
	"net/http"
	"testing"
	"time"
)

func TestHeadTime(t *testing.T) {
	resp, err := http.Head("https://www.time.gov/")
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close() // Always close this without exception

	/* On `Close()` Go implicitly drains any unread bytes.

	The Go HTTP client’s implicit draining of the response body on closing
	could potentially bite you. For example, let’s assume you send a GET request
	for a file and receive a response from the server. You read the response’s
	Content-Length header and realize the file is much larger than you anticipated.
	If you close the response body without reading any of its bytes, Go will
	download the entire file from the server as it drains the body regardless.

	A better alternative would be to send a HEAD request to retrieve the Content-
	Length header. This way, no unread bytes exist in the response body, so closing
	the response body will not incur any additional overhead while draining it.

	On the rare occasion that you make an HTTP request and want to explicitly
	drain the response body, the most effecient way is to use the `io.Copy` function:
	```go
	_, _ = io.Copy(ioutil.Discard, response.Body)
	_ = response.Close()
	```
	The io.Copy function drains the response.Body by reading all bytes from it
	and writing those bytes to ioutil.Discard. As its name indicates, ioutil.Discard
	is a special io.Writer that discards all bytes written to it.
	*/

	now := time.Now().Round(time.Second)
	date := resp.Header.Get("Date")
	if date == "" {
		t.Fatal("no Date header recieved from time.gov")
	}

	dt, err := time.Parse(time.RFC1123, date)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("time.gov: %s (skew %s)", dt, now.Sub(dt))
}
