package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", createSession)
	mux.HandleFunc("/session/", getSession)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server has started on Port: 8080")
	log.Fatal(server.ListenAndServe())
}
