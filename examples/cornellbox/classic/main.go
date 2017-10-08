package main

import (
	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

func scene() Scene {
	return Scene{
		Camera: Cam(1.3),
		Objects: Group(
			ShortBlock(),
			TallBlock(),
			Floor,
			Ceiling,
			BackWall,
			LeftWall.With(ColourRGB(Red)),
			RightWall.With(ColourRGB(Green)),
			CeilingLight().With(Emittance(5.0)),
			Tube(
				Vect(0.1, 0.1, -0.5),
				Vect(0.9, 0.9, -0.5),
				0.2,
			),
		),
	}
}

func main() {
	Run("cornellbox_classic", scene())
}
