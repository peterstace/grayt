package main

import (
	. "github.com/peterstace/grayt/grayt"
)

func main() {

	Run("debug_triangle", Scene{
		Camera: Camera().With(
			Location(Vect(0, 0, 1)),
			LookingAt(Vect(0, 0, 0)),
			FieldOfViewInDegrees(90),
		),
		Objects: Group(
			AlignedSquare(
				Vect(-1, 1, -1),
				Vect(1, -1, -1),
			),
		).With(
			Emittance(1.0),
		),
	})
}
