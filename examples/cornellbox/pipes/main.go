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
			Group(
				Pipe(
					Vect(0.4, 0.1, -0.8),
					Vect(0.5, 0.6, -0.6),
					0.15,
				),
			),
		),
	}
}

func main() {
	Run("cornellbox_pipe", scene())
}
