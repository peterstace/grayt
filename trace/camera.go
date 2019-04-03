package trace

import (
	"math"
	"math/rand"

	"github.com/peterstace/grayt/scene"
	"github.com/peterstace/grayt/xmath"
)

type camera struct {
	screen, eye struct {
		// X vectors go from the center of the screen or eye to the right of
		// the screen or eye.  Y vectors go from the center of the screen or
		// eye towards the top of the screen or eye.  Loc is the location of
		// the center of the screen or eye.
		loc, x, y xmath.Vector
	}
}

func newCamera(conf scene.Camera) camera {
	cam := camera{}

	upDirection := conf.UpDirection.Unit()
	viewDirection := conf.LookingAt.Sub(conf.Location).Unit()

	cam.screen.x = viewDirection.Cross(upDirection)
	cam.screen.y = cam.screen.x.Cross(viewDirection)

	cam.eye.x = cam.screen.x.Scale(conf.FocalLength / conf.FocalRatio)
	cam.eye.y = cam.screen.y.Scale(conf.FocalLength / conf.FocalRatio)
	cam.eye.loc = conf.Location

	halfScreenWidth := math.Tan(conf.FieldOfViewInRadians/2) * conf.FocalLength
	cam.screen.x = cam.screen.x.Scale(halfScreenWidth)
	cam.screen.y = cam.screen.y.Scale(halfScreenWidth)
	cam.screen.loc = cam.eye.loc.Add(viewDirection.Scale(conf.FocalLength))

	return cam
}

func (c *camera) makeRay(x, y float64, rng *rand.Rand) xmath.Ray {
	start := c.eye.loc.
		Add(c.eye.x.Scale(2*rng.Float64() - 1.0)).
		Add(c.eye.y.Scale(2*rng.Float64() - 1.0))
	end := c.screen.loc.
		Add(c.screen.x.Scale(x)).
		Add(c.screen.y.Scale(y))
	return xmath.Ray{
		Start: start,
		Dir:   end.Sub(start),
	}
}
