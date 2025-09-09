package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// How You Write the Response Matters
func TestHandlerWriteHeader(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Bad request"))
		w.WriteHeader(http.StatusBadRequest)
	}
	r := httptest.NewRequest(http.MethodGet, "http://test", nil)
	w := httptest.NewRecorder()
	handler(w, r)
	t.Logf("Response status: %q", w.Result().Status) // 200 OK

	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request"))
	}

	r = httptest.NewRequest(http.MethodGet, "http://test", nil)
	w = httptest.NewRecorder()
	handler(w, r)
	t.Logf("Response status: %q", w.Result().Status) // 400 Bad Request
}
