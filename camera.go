package grayt

import (
	"errors"
	"math"
	"math/rand"
)

type Projection string

const (
	Rectilinear Projection = "Rectilinear"
)

// CameraConfig gives configuration options that are common to all camera
// types. This struct is a parameter to the camera factory functions, primarily
// so in the calling context, it's clear (from the name) whach each
// configuration option is (compared to just passing in Vects and float64s to
// the factory function).
type CameraConfig struct {
	Projection    Projection
	Location      Vector
	ViewDirection Vector
	UpDirection   Vector
	FieldOfView   float64
	FocalLength   float64 // Distance to the focus plane.
	FocalRatio    float64
}

// Camera implementations produce rays from their eye location through a
// virtual screen. Each implementation may have different screen geometry.
type Camera interface {

	// MakeRay produces a ray that goes from the eye to a point on the virtual
	// screen. The left side of the screen has x coordinate -1 and the right
	// side of the screen has x coordinate +1. The top and bottom of the screen
	// have +v and -v respectively, where the value of v depends on the aspect
	// ratio of the screen.
	MakeRay(x, y float64) Ray
}

func NewCamera(conf CameraConfig) (Camera, error) {
	switch conf.Projection {
	case Rectilinear:
		return newRectilinearCamera(conf), nil
	default:
		return nil, errors.New("unknown projection: " + string(conf.Projection))
	}
}

type rectCamera struct {
	screen, eye struct {
		// X vectors go from the center of the screen or eye to the right of
		// the screen or eye.  Y vectors go from the center of the screen or
		// eye towards the top of the screen or eye.  Loc is the location of
		// the center of the screen or eye.
		loc, x, y Vector
	}
}

// NewRectilinearCamera creates a rectilinear camera from a camera config.
func newRectilinearCamera(conf CameraConfig) Camera {

	cam := &rectCamera{}

	conf.UpDirection = conf.UpDirection.Unit()
	conf.ViewDirection = conf.ViewDirection.Unit()

	cam.screen.x = conf.ViewDirection.Cross(conf.UpDirection)
	cam.screen.y = cam.screen.x.Cross(conf.ViewDirection)

	cam.eye.x = cam.screen.x.Scale(conf.FocalLength / conf.FocalRatio)
	cam.eye.y = cam.screen.y.Scale(conf.FocalLength / conf.FocalRatio)
	cam.eye.loc = conf.Location

	halfScreenWidth := math.Tan(conf.FieldOfView/2) * conf.FocalLength
	cam.screen.x = cam.screen.x.Scale(halfScreenWidth)
	cam.screen.y = cam.screen.y.Scale(halfScreenWidth)
	cam.screen.loc = cam.eye.loc.Add(conf.ViewDirection.Scale(conf.FocalLength))

	return cam
}

func (c *rectCamera) MakeRay(x, y float64) Ray {
	start := c.eye.loc.
		Add(c.eye.x.Scale(2*rand.Float64() - 1.0)).
		Add(c.eye.y.Scale(2*rand.Float64() - 1.0))
	end := c.screen.loc.
		Add(c.screen.x.Scale(x)).
		Add(c.screen.y.Scale(y))
	return Ray{
		Start: start,
		Dir:   end.Sub(start),
	}
}
