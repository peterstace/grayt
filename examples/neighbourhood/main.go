package main

import (
	"math/rand"
	"sort"

	. "github.com/peterstace/grayt/grayt"
)

var focus = Vect(-0.3099026210985629, 0.12024389688673498, 0.48822368691539886)

func scene() Scene {
	platform := AlignedBox(
		Vect(-3, -1, -6),
		Vect(+1, 0, +2),
	)

	placeholder := Sphere(
		Vect(0, 0.5, 0),
		0.5,
	)
	placeholder = AlignedBox(
		Vect(-0.5, 0, -0.5),
		Vect(+0.5, 1, +0.5),
	)
	pts := points()
	sort.Slice(pts, func(i, j int) bool { return focus.Sub(pts[i]).LengthSq() < focus.Sub(pts[j]).LengthSq() })
	placeholder = Group(
		balls(pts, pts[0]),
		edges(pts, pts[0]),
	)

	return Scene{
		Camera: Camera().With(
			Location(Vect(3, 5, 15)),
			LookingAt(focus),
			FieldOfViewInDegrees(0.25),
		),
		Objects: Group(
			placeholder,
			platform,
		),
		Sky: Sky(Colour{0.05, 0.05, 0.05}, Colour{4, 4, 4}, Vect(4, 6, 1), 15),
	}
}

func edges(pts []Vector, f Vector) ObjectList {
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
			if !closer && (pts[i] == f || pts[j] == f) {
				edges = Group(edges, Pipe(pts[i], pts[j], 0.01))
			}
		}
	}
	return edges
}

func balls(pts []Vector, f Vector) ObjectList {
	// TODO: Just select one ball
	var objs ObjectList
	for _, p := range pts {
		if p == f {
			objs = Group(objs, Sphere(p, 0.01))
		}
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
