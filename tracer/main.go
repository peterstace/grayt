package main

import (
	"log"
	"os"

	"github.com/peterstace/grayt/scene"
)

func main() {

	f, err := os.Open("foo.bin")
	if err != nil {
		log.Fatal(err)
	}

	s, err := scene.ReadFrom(f)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(s)
}
