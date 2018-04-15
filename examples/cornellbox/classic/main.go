package classic

import (
	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

var SkyFn func(Vector) Colour

func CameraFn() CameraBlueprint {
	return Cam(1.3)
}

func ObjectsFn() ObjectList {
	return Group(
		ShortBlock(),
		TallBlock(),
		Floor,
		Ceiling,
		BackWall,
		LeftWall.With(ColourRGB(Red)),
		RightWall.With(ColourRGB(Green)),
		CeilingLight().With(Emittance(5.0)),
	)
}
