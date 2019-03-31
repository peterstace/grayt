package main

import (
	"log"
	"net/http"
	"os"

	"github.com/peterstace/grayt/api"
	"github.com/peterstace/grayt/mware"
)

const (
	listenAddrEnv = "LISTEN_ADDR"
	assetsDirEnv  = "ASSETS_DIR"
)

func main() {
	listenAddr := os.Getenv(listenAddrEnv)
	if listenAddr == "" {
		log.Fatalf("%s not set", listenAddrEnv)
	}

	assetsDir := os.Getenv(assetsDirEnv)
	if assetsDir == "" {
		log.Fatalf("%s not set", assetsDirEnv)
	}

	s := api.NewServer(assetsDir)
	log.Printf("serving on %v", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mware.LogRequests(s)))
}
