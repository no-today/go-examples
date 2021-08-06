package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server.handler); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
