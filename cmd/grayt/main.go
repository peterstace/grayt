package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/peterstace/grayt/examples/cornellbox/classic"
	"github.com/peterstace/grayt/examples/cornellbox/reflections"
	"github.com/peterstace/grayt/examples/cornellbox/spheretree"
	"github.com/peterstace/grayt/examples/cornellbox/splitbox"
	"github.com/peterstace/grayt/grayt"
)

func main() {
	httpAddr := flag.String("h", ":8080", "http address to listen on")
	storageDir := flag.String("d", "data", "storage directory")
	flag.Parse()

	if err := os.Mkdir(*storageDir, 0751); err != nil && !os.IsExist(err) {
		log.Fatalf("creating storage dir: %v", err)
	}

	s := grayt.NewServer()

	s.Register("cornellbox_classic", classic.CameraFn(), classic.ObjectsFn)
	s.Register("cornellbox_reflections", reflections.CameraFn(), reflections.ObjectsFn)
	s.Register("spheretree", spheretree.CameraFn(), spheretree.ObjectsFn)
	s.Register("splitbox", splitbox.CameraFn(), splitbox.ObjectsFn)

	log.Println("loading...")
	if err := s.Load(*storageDir); err != nil {
		log.Fatalf("could not load server: %v", err)
	}
	log.Println("done")
	go func() {
		for {
			time.Sleep(time.Minute)
			log.Println("saving...")
			if err := s.Save(*storageDir); err != nil {
				log.Fatalf("could not save server: %v", err)
			}
			log.Println("done")
		}
	}()

	if err := s.ListenAndServe(*httpAddr); err != nil {
		log.Fatal(err)
	}
}
