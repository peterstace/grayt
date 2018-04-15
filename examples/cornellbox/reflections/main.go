package reflections

import (
	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

const (
	d = 1.3
	e = 0.05
)

var SkyFn func(Vector) Colour

func CameraFn() CameraBlueprint {
	return Cam(d)
}

func ObjectsFn() ObjectList {
	return Group(

		// Ceiling
		AlignedSquare(Vect(0, 1, d), Vect(1, 1, -1)),
		CeilingLight().With(Emittance(5.0)),

		// Floor
		AlignedSquare(Vect(0, 0, d), Vect(1, 0, -1)),

		// Back wall
		AlignedSquare(Vect(e, e, -1), Vect(1-e, 1-e, -1)).With(Mirror()),
		AlignedSquare(Vect(0, 0, -1), Vect(e, 1, -1)),
		AlignedSquare(Vect(1, 0, -1), Vect(1-e, 1, -1)),
		AlignedSquare(Vect(e, 1-e, -1), Vect(1-e, 1, -1)),
		AlignedSquare(Vect(e, 0, -1), Vect(1-e, e, -1)),

		// Front wall
		AlignedSquare(Vect(0, 0, d), Vect(1, 1, d)),

		// Walls

		Group(
			AlignedSquare(Vect(0, 0, 0), Vect(0, 1, -e)),
			AlignedSquare(Vect(0, 0, -1+e), Vect(0, 1, -1)),
			AlignedSquare(Vect(0, 0, -e), Vect(0, e, -1+e)),
			AlignedSquare(Vect(0, 1-e, -e), Vect(0, 1, -1+e)),
		).With(ColourRGB(Red)),
		AlignedSquare(Vect(0, e, -e), Vect(0, 1-e, -1+e)).With(Mirror()),

		Group(
			AlignedSquare(Vect(1, 0, 0), Vect(1, 1, -e)),
			AlignedSquare(Vect(1, 0, -1+e), Vect(1, 1, -1)),
			AlignedSquare(Vect(1, 0, -e), Vect(1, e, -1+e)),
			AlignedSquare(Vect(1, 1-e, -e), Vect(1, 1, -1+e)),
		).With(ColourRGB(Green)),
		AlignedSquare(Vect(1, e, -e), Vect(1, 1-e, -1+e)).With(Mirror()),

		AlignedSquare(Vect(0, 0, 0), Vect(0, 1, d)),
		AlignedSquare(Vect(1, 0, 0), Vect(1, 1, d)),

		// Blocks
		ShortBlock().With(ColourRGB(0xd684ff)),
		TallBlock().With(ColourRGB(0x3ae3ff)),

		// Spheres
		Group(
			Sphere(Vect(0.2, 0.1, -0.2), 0.1),
			Sphere(Vect(0.85, 0.1, -0.85), 0.1),
			Sphere(Vect(0.33, 0.6+0.075, -0.64), 0.075),
		).With(Mirror()),
	)
}
