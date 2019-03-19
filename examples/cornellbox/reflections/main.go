package reflections

import (
	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
	"github.com/peterstace/grayt/xmath"
)

const (
	d = 1.3
	e = 0.05
)

func CameraFn() CameraBlueprint {
	return Cam(d)
}

func ObjectsFn() ObjectList {
	return Group(

		// Ceiling
		AlignedSquare(xmath.Vect(0, 1, d), xmath.Vect(1, 1, -1)),
		CeilingLight().With(Emittance(5.0)),

		// Floor
		AlignedSquare(xmath.Vect(0, 0, d), xmath.Vect(1, 0, -1)),

		// Back wall
		AlignedSquare(xmath.Vect(e, e, -1), xmath.Vect(1-e, 1-e, -1)).With(Mirror()),
		AlignedSquare(xmath.Vect(0, 0, -1), xmath.Vect(e, 1, -1)),
		AlignedSquare(xmath.Vect(1, 0, -1), xmath.Vect(1-e, 1, -1)),
		AlignedSquare(xmath.Vect(e, 1-e, -1), xmath.Vect(1-e, 1, -1)),
		AlignedSquare(xmath.Vect(e, 0, -1), xmath.Vect(1-e, e, -1)),

		// Front wall
		AlignedSquare(xmath.Vect(0, 0, d), xmath.Vect(1, 1, d)),

		// Walls

		Group(
			AlignedSquare(xmath.Vect(0, 0, 0), xmath.Vect(0, 1, -e)),
			AlignedSquare(xmath.Vect(0, 0, -1+e), xmath.Vect(0, 1, -1)),
			AlignedSquare(xmath.Vect(0, 0, -e), xmath.Vect(0, e, -1+e)),
			AlignedSquare(xmath.Vect(0, 1-e, -e), xmath.Vect(0, 1, -1+e)),
		).With(ColourRGB(Red)),
		AlignedSquare(xmath.Vect(0, e, -e), xmath.Vect(0, 1-e, -1+e)).With(Mirror()),

		Group(
			AlignedSquare(xmath.Vect(1, 0, 0), xmath.Vect(1, 1, -e)),
			AlignedSquare(xmath.Vect(1, 0, -1+e), xmath.Vect(1, 1, -1)),
			AlignedSquare(xmath.Vect(1, 0, -e), xmath.Vect(1, e, -1+e)),
			AlignedSquare(xmath.Vect(1, 1-e, -e), xmath.Vect(1, 1, -1+e)),
		).With(ColourRGB(Green)),
		AlignedSquare(xmath.Vect(1, e, -e), xmath.Vect(1, 1-e, -1+e)).With(Mirror()),

		AlignedSquare(xmath.Vect(0, 0, 0), xmath.Vect(0, 1, d)),
		AlignedSquare(xmath.Vect(1, 0, 0), xmath.Vect(1, 1, d)),

		// Blocks
		ShortBlock().With(ColourRGB(0xd684ff)),
		TallBlock().With(ColourRGB(0x3ae3ff)),

		// Spheres
		Group(
			Sphere(xmath.Vect(0.2, 0.1, -0.2), 0.1),
			Sphere(xmath.Vect(0.85, 0.1, -0.85), 0.1),
			Sphere(xmath.Vect(0.33, 0.6+0.075, -0.64), 0.075),
		).With(Mirror()),
	)
}
