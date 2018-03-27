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

	s.Register("cornellbox_classic", classic.SceneFn)
	s.Register("cornellbox_reflections", reflections.SceneFn)
	s.Register("spheretree", spheretree.SceneFn)
	s.Register("splitbox", splitbox.SceneFn)
	s.Register("neighbourhood", neighbourhood.SceneFn)

	if err := s.ListenAndServe(*httpAddr); err != nil {
		log.Fatal(err)
	}
}
