package classic

import (
	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

func SceneFn() Scene {
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
		),
	}
}
