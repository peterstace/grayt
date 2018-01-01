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

	grayt.Register("cornellbox_classic", classic.SceneFn)
	grayt.Register("cornellbox_reflections", reflections.SceneFn)
	grayt.Register("spheretree", spheretree.SceneFn)
	grayt.Register("splitbox", splitbox.SceneFn)
	grayt.Register("neighbourhood", neighbourhood.SceneFn)

	if err := grayt.ListenAndServe(*httpAddr); err != nil {
		log.Fatal(err)
	}
}
