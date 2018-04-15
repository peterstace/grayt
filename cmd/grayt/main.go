package main

import (
	"flag"
	"log"

	"github.com/peterstace/grayt/examples/cornellbox/classic"
	"github.com/peterstace/grayt/examples/cornellbox/reflections"
	"github.com/peterstace/grayt/examples/cornellbox/spheretree"
	"github.com/peterstace/grayt/examples/cornellbox/splitbox"
	"github.com/peterstace/grayt/examples/neighbourhood"
	"github.com/peterstace/grayt/grayt"
)

func main() {
	httpAddr := flag.String("h", ":8080", "http address to listen on")
	flag.Parse()

	s := grayt.NewServer()

	s.Register("cornellbox_classic", classic.SkyFn, classic.CameraFn(), classic.ObjectsFn)
	s.Register("cornellbox_reflections", reflections.SkyFn, reflections.CameraFn(), reflections.ObjectsFn)
	s.Register("spheretree", spheretree.SkyFn, spheretree.CameraFn(), spheretree.ObjectsFn)
	s.Register("splitbox", splitbox.SkyFn, splitbox.CameraFn(), splitbox.ObjectsFn)
	s.Register("neighbourhood", neighbourhood.SkyFn, neighbourhood.CameraFn(), neighbourhood.ObjectsFn)

	if err := s.ListenAndServe(*httpAddr); err != nil {
		log.Fatal(err)
	}
}
