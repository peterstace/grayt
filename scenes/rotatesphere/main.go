package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/peterstace/grayt"
)

func main() {
	i := flag.Int("i", 0, "index")
	n := flag.Int("n", 60, "range")
	flag.Parse()

	r := grayt.NewRunner()
	r.PxWide = 200
	r.PxHigh = r.PxWide * 3 / 4
	r.Quality = 1 << 10
	r.BaseName = fmt.Sprintf("RotateSphere_%03d_%03d", *i, *n)
	r.Run(scene(float64(*i) / float64(*n)))
}

/*

	+----
	|
	-3 -2  -1   0   1   2   3
	|   *   *   *   *   *   |
	|
	|
	+

*/

func scene(progress float64) grayt.Scene {

	/*
		+----------------+---+
		|                |   |
		+---+------------+   |
		|   |            |   |
		|   |            |   |
		|   | x x x      |   |
		|   |            |   |
		|   |            |   |
		|   +------------+---+
		|   |                |
		+---+----------------+
	*/

	progress *= 2.0 * math.Pi

	surfaces := []grayt.Surface{
		grayt.AlignedBox(
			grayt.Vect(-4, 0, -3),
			grayt.Vect(3, 1, -4),
		),
		grayt.AlignedBox(
			grayt.Vect(-3, 0, 3),
			grayt.Vect(4, 1, 4),
		),
		grayt.AlignedBox(
			grayt.Vect(-3, 0, 4),
			grayt.Vect(-4, 1, -3),
		),
		grayt.AlignedBox(
			grayt.Vect(3, 0, -4),
			grayt.Vect(4, 1, 3),
		),
		grayt.YPlane(),
		grayt.Sphere(1).Translate(
			math.Sin(progress-math.Pi/2)-1,
			1,
			-math.Sin(progress),
		),
		grayt.Sphere(1).Translate(
			math.Sin(progress-math.Pi/2)+1,
			1,
			math.Sin(progress),
		),
	}

	var entities []grayt.Entity
	for _, surface := range surfaces {
		entities = append(entities, grayt.Entity{
			Surface:  surface,
			Material: grayt.Material{Colour: grayt.White},
		})
	}
	entities = append(entities, grayt.Entity{
		Surface:  grayt.Sphere(3).Translate(0, 10, 0),
		Material: grayt.Material{Colour: grayt.White, Emittance: 5},
	})

	return grayt.Scene{Camera: cam(), Entities: entities}
}

func cam() grayt.Camera {
	loc := grayt.Vect(-12, 16, 20)
	return grayt.NewRectilinearCamera(grayt.CameraConfig{
		Location:      loc,
		ViewDirection: grayt.Vect(0.0, 0.0, 0.0).Sub(loc),
		UpDirection:   grayt.Vect(0, 1, 0),
		FieldOfView:   0.55,
		FocalLength:   10,
		FocalRatio:    math.Inf(+1),
	})
}
