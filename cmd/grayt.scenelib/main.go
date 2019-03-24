package main

import (
	"log"
	"net/http"
	"os"

	"github.com/peterstace/grayt/mware"
	"github.com/peterstace/grayt/scenelib"
	"github.com/peterstace/grayt/scenelib/cornellbox"
)

const (
	listenAddrEnv = "GRAYT_SCENELIB_LISTEN_ADDR"
)

func main() {
	listenAddr, ok := os.LookupEnv(listenAddrEnv)
	if !ok {
		log.Fatalf("%s not set", listenAddrEnv)
	}

	s := scenelib.NewServer()
	s.Register("cornellbox_classic", cornellbox.Classic)
	/*
		s.Register("cornellbox_reflections", reflections.CameraFn(), reflections.ObjectsFn)
		s.Register("spheretree", spheretree.CameraFn(), spheretree.ObjectsFn)
		s.Register("splitbox", splitbox.CameraFn(), splitbox.ObjectsFn)
	*/

	log.Printf("serving on %v", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mware.LogRequests(s)))
}
