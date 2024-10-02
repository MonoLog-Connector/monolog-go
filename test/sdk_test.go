package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MonoLog-Connector/monolog-go/client"
)

// A simple test for the TrackRequestTime middleware
func TestTrackRequestTime(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	wrappedHandler := client.TrackRequestTime(handler)
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
