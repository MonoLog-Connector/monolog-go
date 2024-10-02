package client

import (
	"log"
	"net/http"
	"time"
)

// TrackRequestTime is middleware to track the time taken for each HTTP request
func TrackRequestTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the time taken
		duration := time.Since(startTime)
		log.Printf("Request to %s took %v", r.URL.Path, duration)
	})
}
