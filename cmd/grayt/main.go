package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/peterstace/grayt/examples/cornellbox/classic"
	"github.com/peterstace/grayt/examples/cornellbox/reflections"
	"github.com/peterstace/grayt/examples/cornellbox/spheretree"
	"github.com/peterstace/grayt/examples/cornellbox/splitbox"
	"github.com/peterstace/grayt/examples/neighbourhood"
	"github.com/peterstace/grayt/grayt"
)

var (
	pxWide     = flag.Int("w", 640, "width in pixels")
	quality    = flag.Int("q", 10, "quality (samples per pixel)")
	verbose    = flag.Bool("v", false, "verbose model")
	output     = flag.String("o", "", "output file override")
	numWorkers = flag.Int("j", runtime.GOMAXPROCS(0), "number of worker goroutines")
	debug      = flag.Bool("d", false, "debug mode (enable assertions)")
	normals    = flag.Bool("n", false, "plot normals")
	scene      = flag.String("s", "", "scene to render")
	httpAddr   = flag.String("h", ":8080", "http address to listen on")
)

func main() {
	flag.Parse()
	grayt.Config.PxWide = *pxWide
	grayt.Config.Quality = *quality
	grayt.Config.Verbose = *verbose
	grayt.Config.Output = *output
	grayt.Config.NumWorkers = *numWorkers
	grayt.Config.EnableAssertions = *debug
	grayt.Config.Normals = *normals
	grayt.Config.Scene = *scene

	grayt.Register("cornellbox_classic", classic.SceneFn)
	grayt.Register("cornellbox_reflections", reflections.SceneFn)
	grayt.Register("spheretree", spheretree.SceneFn)
	grayt.Register("splitbox", splitbox.SceneFn)
	grayt.Register("neighbourhood", neighbourhood.SceneFn)

	if err := grayt.ListenAndServe(*httpAddr); err != nil {
		log.Fatal(err)
	}
}
