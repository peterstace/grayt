package tracer

import (
	"math"
	"math/rand"

	"github.com/peterstace/grayt/ray"
	"github.com/peterstace/grayt/vect"
)

type rectilinearCamera struct {
	screen, eye struct {
		// X vectors go from the center of the screen or eye to the right of
		// the screen or eye.  Y vectors go from the center of the screen or
		// eye towards the top of the screen or eye.  Loc is the location of
		// the center of the screen or eye.
		loc, x, y vect.V
	}
}

// NewRectilinearCamera creates a rectilinear camera from a camera config.
func NewRectilinearCamera(conf CameraConfig) Camera {

	cam := &rectilinearCamera{}

	conf.UpDirection = conf.UpDirection.Unit()
	conf.ViewDirection = conf.ViewDirection.Unit()

	cam.screen.x = vect.Cross(conf.ViewDirection, conf.UpDirection)
	cam.screen.y = vect.Cross(cam.screen.x, conf.ViewDirection)

	cam.eye.x = cam.screen.x.Extended(conf.FocalLength / conf.FocalRatio)
	cam.eye.y = cam.screen.y.Extended(conf.FocalLength / conf.FocalRatio)
	cam.eye.loc = conf.Location

	halfScreenWidth := math.Tan(conf.FieldOfView/2) * conf.FocalLength
	cam.screen.x = cam.screen.x.Extended(halfScreenWidth)
	cam.screen.y = cam.screen.y.Extended(halfScreenWidth)
	cam.screen.loc = vect.Add(cam.eye.loc, conf.ViewDirection.Extended(conf.FocalLength))

	return cam
}

func (c *rectilinearCamera) MakeRay(x, y float64) ray.Ray {
	start := vect.Add(
		c.eye.loc,
		vect.Add(
			c.eye.x.Extended(2*rand.Float64()-1.0),
			c.eye.y.Extended(2*rand.Float64()-1.0),
		),
	)
	end := vect.Add(
		c.screen.loc,
		vect.Add(
			c.screen.x.Extended(x),
			c.screen.y.Extended(y),
		),
	)
	return ray.Ray{
		Start: start,
		Dir:   vect.Sub(end, start).Unit(),
	}
}
