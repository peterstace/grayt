package main

import (
	"log"
	"net/http"
	"os"

	"github.com/peterstace/grayt/mware"
	"github.com/peterstace/grayt/worker"
)

const (
	listenAddrEnv = "LISTEN_ADDR"
)

func main() {
	listenAddr := os.Getenv(listenAddrEnv)
	if listenAddr == "" {
		log.Fatalf("%s not set", listenAddrEnv)
	}

	s := worker.NewServer()
	log.Printf("serving on %v", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mware.LogRequests(s)))
}
