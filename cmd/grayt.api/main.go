package main

import (
	"log"
	"net/http"
	"os"

	"github.com/peterstace/grayt/api"
	"github.com/peterstace/grayt/mware"
)

const (
	listenAddrEnv   = "LISTEN_ADDR"
	scenelibAddrEnv = "SCENELIB_ADDR"
	workerAddrEnv   = "WORKER_ADDR"
	assetsDirEnv    = "ASSETS_DIR"
)

func main() {
	listenAddr := os.Getenv(listenAddrEnv)
	if listenAddr == "" {
		log.Fatalf("%s not set", listenAddrEnv)
	}

	scenelibAddr := os.Getenv(scenelibAddrEnv)
	if scenelibAddr == "" {
		log.Fatalf("%s not set", scenelibAddrEnv)
	}

	workerAddr := os.Getenv(workerAddrEnv)
	if workerAddr == "" {
		log.Fatalf("%s not set", workerAddrEnv)
	}

	assetsDir := os.Getenv(assetsDirEnv)
	if assetsDir == "" {
		log.Fatalf("%s not set", assetsDirEnv)
	}

	s := api.NewServer(scenelibAddr, workerAddr, assetsDir)
	log.Printf("serving on %v", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mware.LogRequests(s)))
}
