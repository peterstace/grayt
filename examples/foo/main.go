package main

import (
	"log"
	"os"

	"github.com/peterstace/grayt/scene"
)

func main() {
	s := scene.Scene{
		Camera: scene.DefaultCamera(),
		Triangles: []scene.Triangle{
			{
				A:         scene.Vector{X: -0.5, Y: -0.5, Z: -5.0},
				B:         scene.Vector{X: 0.5, Y: -0.5, Z: -5.0},
				C:         scene.Vector{X: 0.5, Y: 0.5, Z: -5.0},
				Colour:    [3]float64{1, 1, 1},
				Emittance: 1.0,
			},
		},
	}

	f, err := os.Create("foo.bin")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.WriteTo(f); err != nil {
		log.Fatal(err)
	}
}
