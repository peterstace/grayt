package main

import (
	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("test", scene())
}

var focus = Vect(-0.3099026210985629, 0.12024389688673498, 0.48822368691539886)

func scene() Scene {
	sp := Sphere(Vect(
		-0.3131815596009845,
		0.11208661069616384,
		0.4834468154653815,
	), 0.01)

	p1 := Pipe(
		Vect(
			-0.2466522683500979,
			0.056480273967004564,
			0.4429233173416114,
		),
		Vect(
			-0.3131815596009845,
			0.11208661069616384,
			0.4834468154653815,
		), 0.01)

	_ = sp
	_ = p1

	return Scene{
		Camera: Camera().With(
			Location(Vect(3, 5, 15)),
			LookingAt(focus),
			FieldOfViewInDegrees(0.1),
		),
		Objects: Group(
			p1,
			sp,
		),
		Sky: func(Vector) Colour { return Colour{0.05, 0.05, 0.05} },
	}
}
