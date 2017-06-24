package main

import (
	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

func scene() Scene {
	return Scene{
		Camera: Cam(1.3),
		Objects: Group(
			ShortBlock().With(Translate(Vect(-0.5, 0, 0))),
			TallBlock().With(Translate(Vect(0.5, 0, 0))),
			Floor,
			Ceiling,
			BackWall,
			LeftWall.With(ColourRGB(Red)),
			RightWall.With(ColourRGB(Green)),
			CeilingLight().With(Emittance(5.0)),
		),
	}
}

func main() {
	Run("cornellbox_classic_shift", scene())
}
