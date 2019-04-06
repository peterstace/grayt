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
		log.Fatal("LISTEN_ADDR not set")
	}

	assetsDir := os.Getenv("ASSETS_DIR")
	if assetsDir == "" {
		log.Fatal("ASSETS_DIR not set")
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		log.Fatal("DATA_DIR not set")
	}

	s := api.NewServer(assetsDir, dataDir)
	log.Printf("serving on %v", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, s))
}
