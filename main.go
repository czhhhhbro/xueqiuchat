package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

type Message struct {
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	Self      bool      `json:"self"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	messages []Message
	mu       sync.Mutex
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		msg.Timestamp = time.Now()
		msg.Self = true

		mu.Lock()
		messages = append(messages, msg)
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		list := messages
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
