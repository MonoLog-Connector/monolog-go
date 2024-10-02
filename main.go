package main

import (
	"log"
	"net/http"

	"github.com/MonoLog-Connector/monolog-go/client"
)

func main() {
	mux := http.NewServeMux()
	// Example handler
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	// Wrap the handler with the middleware to track request time
	wrappedMux := client.TrackRequestTime(mux)
	// Start the server
	log.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", wrappedMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
