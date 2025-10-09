package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	IP        string    `json:"ip"`
	Port      string    `json:"port"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Checksum  string    `json:"checksum"`
}

type RegisterSessionResp struct {
	Message string        `json:"message"`
	Code    string        `json:"code"`
	Timeout time.Duration `json:"timeout"`
}

type QueryResp struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Checksum string `json:"checksum"`
}

var (
	sessions = make(map[string]Session)
	mut      sync.Mutex
)

func createSession(w http.ResponseWriter, r *http.Request) {
	var s Session
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mut.Lock()
	sessions[s.Code] = s
	mut.Unlock()

	resp := RegisterSessionResp{
		Code:    s.Code,
		Timeout: time.Until(s.ExpiresAt),
		Message: "Session created successfully",
	}

	data, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getSession(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/session/"):]

	log.Println("Code: ", code)

	mut.Lock()
	s, ok := sessions[code]

	if ok {
		// remove that session from server
		delete(sessions, code)
	}
	mut.Unlock()

	if !ok {
		http.Error(w, "Invalid Code", http.StatusBadRequest)
		return
	}

	resp := QueryResp{
		IP:       s.IP,
		Port:     s.Port,
		FileName: s.FileName,
		FileSize: s.FileSize,
		Checksum: s.Checksum,
	}

	data, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
