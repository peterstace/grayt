package main

import (
	"log"
	"net/http"
	"os"

	"github.com/peterstace/grayt/api"
)

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		log.Fatalf("LISTEN_ADDR not set")
	}

	s := api.NewServer()
	log.Printf("serving on %v", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, s))
}
