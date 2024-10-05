package client

import (
	"log"
	"net/http"
	"time"
)

func TrackerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		LogRequestDetails(r, elapsed)
	})
}

func LogRequestDetails(r *http.Request, elapsed time.Duration) {
	log.Printf("Request: %s %s | Elapsed Time: %s\n",
		r.Method, r.URL.Path, elapsed)
}
