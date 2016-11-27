package grayt

import (
	"math"
	"math/rand"
)

type camera struct {
	screen, eye struct {
		// X vectors go from the center of the screen or eye to the right of
		// the screen or eye.  Y vectors go from the center of the screen or
		// eye towards the top of the screen or eye.  Loc is the location of
		// the center of the screen or eye.
		loc, x, y Vector
	}
}

func newCamera(conf CameraBlueprint) camera {

	cam := camera{}

	upDirection := conf.upDirection.Unit()
	viewDirection := conf.lookingAt.Sub(conf.location).Unit()

	cam.screen.x = viewDirection.cross(upDirection)
	cam.screen.y = cam.screen.x.cross(viewDirection)

	cam.eye.x = cam.screen.x.Scale(conf.focalLength / conf.focalRatio)
	cam.eye.y = cam.screen.y.Scale(conf.focalLength / conf.focalRatio)
	cam.eye.loc = conf.location

	halfScreenWidth := math.Tan(conf.fieldOfViewInRadians/2) * conf.focalLength
	cam.screen.x = cam.screen.x.Scale(halfScreenWidth)
	cam.screen.y = cam.screen.y.Scale(halfScreenWidth)
	cam.screen.loc = cam.eye.loc.Add(viewDirection.Scale(conf.focalLength))

	return cam
}

func (c *camera) makeRay(x, y float64) ray {
	start := c.eye.loc.
		Add(c.eye.x.Scale(2*rand.Float64() - 1.0)).
		Add(c.eye.y.Scale(2*rand.Float64() - 1.0))
	end := c.screen.loc.
		Add(c.screen.x.Scale(x)).
		Add(c.screen.y.Scale(y))
	return ray{
		start: start,
		dir:   end.Sub(start),
	}
}
