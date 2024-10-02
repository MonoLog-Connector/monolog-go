package client

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// TrackRequestTime is middleware to track the time taken for each HTTP request
func TrackRequestTime(next http.Handler) http.Handler {
	log.Formatter = &logrus.JSONFormatter{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// Call the next handler
		next.ServeHTTP(w, r)
		// Log the time taken
		duration := time.Since(startTime)
		log.WithFields(logrus.Fields{
			"authKey":   "",
			"RPM":       100,
			"CPU Usage": "50%",
			"url":       r.URL.Path,
			"duration":  duration,
		}).Info("Request tracked")
		// log.Printf("Request to %s took %v", r.URL.Path, duration)
	})
}
