package main

import (
	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

func scene() Scene {
	return Scene{
		Camera: Cam(1.3),
		Objects: Group(
			Floor,
			Ceiling,
			BackWall,
			LeftWall.With(ColourRGB(Red)),
			RightWall.With(ColourRGB(Green)),
			CeilingLight().With(Emittance(5.0)),
			Pipe(
				Vect(0.4, 0.0, -0.5),
				Vect(0.4, 0.3, -0.5),
				0.35,
			),
			Sphere(
				Vect(0.4, 0.3, -0.5),
				0.35,
			),
		),
	}
}

func main() {
	Run("cornellbox_pipe", scene())
}
