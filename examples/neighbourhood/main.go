package main

import (
	"math/rand"

	. "github.com/peterstace/grayt/grayt"
)

func scene() Scene {
	platform := AlignedBox(
		Vect(-3, -1, -6),
		Vect(+1, 0, +2),
	)

	pts := points()
	structure := Group(
		balls(pts),
		edges(pts),
	).With(ColourRGB(0x00aaaaaa))

	return Scene{
		Camera: Camera().With(
			Location(Vect(3, 5, 15)),
			LookingAt(Vect(0, 0.5, 0)),
			FieldOfViewInDegrees(6),
		),
		Objects: Group(
			structure,
			platform,
		),
		Sky: Sky(Colour{0.05, 0.05, 0.05}, Colour{4, 4, 4}, Vect(4, 6, 1), 15),
	}
}

func edges(pts []Vector) ObjectList {
	type edge struct {
		a Vector
		b Vector
	}
	var edges ObjectList
	for i := range pts {
		for j := 0; j < i; j++ {
			ijDist := pts[i].Sub(pts[j]).LengthSq()
			closer := false
			for k := range pts {
				ikDist := pts[i].Sub(pts[k]).LengthSq()
				jkDist := pts[j].Sub(pts[k]).LengthSq()
				if ikDist < ijDist && jkDist < ijDist {
					closer = true
					break
				}
			}
			if !closer {
				edges = Group(edges, Pipe(pts[i], pts[j], 0.01))
			}
		}
	}
	return edges
}

func balls(pts []Vector) ObjectList {
	var objs ObjectList
	for _, p := range pts {
		objs = Group(objs, Sphere(p, 0.01))
	}
	return objs
}

func points() []Vector {
	rnd := rand.New(rand.NewSource(1))
	var pp []Vector
	for i := 0; i < 1000; i++ {
		pp = append(pp, Vector{
			rnd.Float64() - 0.5,
			rnd.Float64(),
			rnd.Float64() - 0.5,
		})
	}
	return pp
}

func main() {
	Run("neighbourhood", scene())
}
