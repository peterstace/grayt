package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/peterstace/grayt/grayt"
)

func main() {
	// TODO: Remove this hack to sleep so that everything can start up.
	time.Sleep(time.Second)

	httpAddr := flag.String("h", ":8080", "http address to listen on")
	storageDir := flag.String("d", "data", "storage directory")
	flag.Parse()

	if err := os.Mkdir(*storageDir, 0751); err != nil && !os.IsExist(err) {
		log.Fatalf("creating storage dir: %v", err)
	}

	s := grayt.NewServer()

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