package scenes

import (
	"math"

	"github.com/peterstace/grayt/engine"
)

func init() {
	engine.RegisterScene("cornell_classic", func(a *engine.API) {
		a.SetAspectRatio(1, 1)

		d := 1.3
		a.CameraLocation(0.5, 0.5, d)
		a.CameraLookingAt(0.5, 0.5, -0.5)
		a.CameraUpDirection(0, 1, 0)
		a.CameraRadFieldOfView(2 * math.Asin(0.5/math.Sqrt(0.25+d*d)))

		a.Illuminate(1)
		a.Tri(
			0.5, 0.9, -0.9,
			0.1, 0.9, -0.1,
			0.9, 0.9, -0.1,
		)
		a.Illuminate(0)
		a.Tri(
			0.5, 0.1, -0.9,
			0.1, 0.1, -0.1,
			0.9, 0.1, -0.1,
		)
	})
}
