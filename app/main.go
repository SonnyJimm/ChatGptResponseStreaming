package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Event struct {
	Data string `json:"data"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Set up a timer to send events every second
	ticker := time.NewTicker(1 * time.Second)
	exit := time.After(10 * time.Second)
	// Send events until the context is canceled

	for {
		select {
		case <-ticker.C:
			event := Event{Data: fmt.Sprintf("The time is now %s", time.Now().Format(time.RFC3339))}
			eventJSON, err := json.Marshal(event)
			if err != nil {
				log.Println("Error marshaling event to JSON:", err)
				continue
			}
			fmt.Fprintf(w, string(eventJSON))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

		case <-exit:
			ticker.Stop()
			return

		case <-r.Context().Done():
			ticker.Stop()
			return
		}
	}
}

func main() {
	// handle every request to sending /
	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Error starting SSE server:", err)
	}
}
