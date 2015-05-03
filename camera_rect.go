package grayt

import (
	"math"
	"math/rand"
)

type rectilinearCamera struct {
	screen, eye struct {
		// X vectors go from the center of the screen or eye to the right of
		// the screen or eye.  Y vectors go from the center of the screen or
		// eye towards the top of the screen or eye.  Loc is the location of
		// the center of the screen or eye.
		loc, x, y Vect
	}
}

// NewRectilinearCamera creates a rectilinear camera from a camera config.
func NewRectilinearCamera(conf CameraConfig) Camera {

	cam := &rectilinearCamera{}

	conf.UpDirection = conf.UpDirection.Unit()
	conf.ViewDirection = conf.ViewDirection.Unit()

	cam.screen.x = conf.ViewDirection.Cross(conf.UpDirection)
	cam.screen.y = cam.screen.x.Cross(conf.ViewDirection)

	cam.eye.x = cam.screen.x.Extended(conf.FocalLength / conf.FocalRatio)
	cam.eye.y = cam.screen.y.Extended(conf.FocalLength / conf.FocalRatio)
	cam.eye.loc = conf.Location

	halfScreenWidth := math.Tan(conf.FieldOfView/2) * conf.FocalLength
	cam.screen.x = cam.screen.x.Extended(halfScreenWidth)
	cam.screen.y = cam.screen.y.Extended(halfScreenWidth)
	cam.screen.loc = cam.eye.loc.Add(conf.ViewDirection.Extended(conf.FocalLength))

	return cam
}

func (c *rectilinearCamera) MakeRay(x, y float64) Ray {
	start := c.eye.loc.
		Add(c.eye.x.Extended(2*rand.Float64() - 1.0)).
		Add(c.eye.y.Extended(2*rand.Float64() - 1.0))
	end := c.screen.loc.
		Add(c.screen.x.Extended(x)).
		Add(c.screen.y.Extended(y))
	return Ray{
		Start: start,
		Dir:   end.Sub(start),
	}
}
