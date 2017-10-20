package main

import (
	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("test", scene())
}

func scene() Scene {

	return Scene{
		Camera: Camera().With(
			Location(Vect(0, 3, 6)),
			LookingAt(Vector{}),
			FieldOfViewInDegrees(30),
		),
		Objects: Group(
			Sphere(Vect(0, 0, 0), 0.1),
			Pipe(Vect(-1, -1, -1), Vect(0, 0, 0), 0.1),
		),
		Sky: func(Vector) Colour { return Colour{1, 1, 1} },
	}
}
