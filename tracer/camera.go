package main

import (
	"math"
	"math/rand"

	"github.com/peterstace/grayt/scene"
)

type camera struct {
	screen, eye struct {
		// X vectors go from the center of the screen or eye to the right of
		// the screen or eye.  Y vectors go from the center of the screen or
		// eye towards the top of the screen or eye.  Loc is the location of
		// the center of the screen or eye.
		loc, x, y vector
	}
}

func newCamera(conf scene.Camera) camera {

	cam := camera{}

	upDirection := convertVector(conf.UpDirection).unit()
	viewDirection := convertVector(conf.ViewDirection).unit()

	cam.screen.x = viewDirection.cross(upDirection)
	cam.screen.y = cam.screen.x.cross(viewDirection)

	cam.eye.x = cam.screen.x.scale(conf.FocalLength / conf.FocalRatio)
	cam.eye.y = cam.screen.y.scale(conf.FocalLength / conf.FocalRatio)
	cam.eye.loc = convertVector(conf.Location)

	halfScreenWidth := math.Tan(conf.FieldOfViewInDegrees*math.Pi/180/2) * conf.FocalLength
	cam.screen.x = cam.screen.x.scale(halfScreenWidth)
	cam.screen.y = cam.screen.y.scale(halfScreenWidth)
	cam.screen.loc = cam.eye.loc.add(viewDirection.scale(conf.FocalLength))

	return cam
}

func (c *camera) makeRay(x, y float64) ray {
	start := c.eye.loc.
		add(c.eye.x.scale(2*rand.Float64() - 1.0)).
		add(c.eye.y.scale(2*rand.Float64() - 1.0))
	end := c.screen.loc.
		add(c.screen.x.scale(x)).
		add(c.screen.y.scale(y))
	return ray{
		start: start,
		dir:   end.sub(start),
	}
}
