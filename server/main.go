package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

type Key string

type Session struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	Key  *Key   `json:"key"` // public key for encrypt or decrypt the file (currently lets make it nil)
}

var (
	sessions = make(map[string]Session)
	mut      sync.Mutex
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/session", createSession)
	mux.HandleFunc("/session/", getSession)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server has started on Port: 8080")
	log.Fatal(server.ListenAndServe())
}

func generateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func createSession(w http.ResponseWriter, r *http.Request) {
	// create a session on server and return the 6-digit code
	var s Session
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code := generateCode()

	mut.Lock() // possibilty that multiple user trying to access sessions
	sessions[code] = s
	mut.Unlock()

	json.NewEncoder(w).Encode(map[string]string{"code": code})
}

func getSession(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/session/"):]

	mut.Lock()
	s, ok := sessions[code]

	if ok {
		// remove that session from server
		delete(sessions, code)
	}
	mut.Unlock()

	if !ok {
		http.Error(w, "Code not found or Invalid Code", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(s)
}
