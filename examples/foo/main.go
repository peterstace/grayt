package main

import (
	"log"
	"os"

	"github.com/peterstace/grayt/scene"
)

func main() {
	s := scene.Scene{
		Camera: scene.DefaultCamera(),
		Triangles: []scene.Triangle{{
			A:         scene.Vector{},
			B:         scene.Vector{},
			C:         scene.Vector{},
			Colour:    [3]uint16{0, 0, 0},
			Emittance: 1.0,
		}},
	}

	f, err := os.Create("foo.bin")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.WriteTo(f); err != nil {
		log.Fatal(err)
	}
}
