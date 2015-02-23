package tracer

import (
	"math"
	"math/rand"

	"github.com/peterstace/grayt/ray"
	"github.com/peterstace/grayt/vect"
)

// Camera implementations produce rays that go from an eye to a virtual screen.
// The rays produced are specified via a coordiate system on the virtual
// screen.  The left side of the virtual screen has x coordinate -1, the right
// side of the virtual screen has coordinate +1. The top of the virtual screen
// has y coordinate v and the bottom of the virtual screen has y coordinate -v
// (where the value of v depends on the aspect ratio of the screen).
type Camera interface {
	MakeRay(x, y float64) ray.Ray
}

type rectilinearCamera struct {
	screen, eye struct {
		// X vectors go from the center of the screen or eye to the right of
		// the screen or eye.  Y vectors go from the center of the screen or
		// eye towards the top of the screen or eye.  Loc is the location of
		// the center of the screen or eye.
		loc, x, y vect.V
	}
}

// CameraConfig gives configuration options that are common to all camera
// types. This struct is a parameter to the camera factory functions, primarily
// so in the calling context, it's clear (from the name) whach each
// configuration option is (compared to just passing in Vects and float64s to
// the factory function).
type CameraConfig struct {
	Location      vect.V
	ViewDirection vect.V
	UpDirection   vect.V
	FieldOfView   float64
	FocalLength   float64
	FocalRatio    float64
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
