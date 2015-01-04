package grayt

import "math"

// Cameras produce rays that go from an eye to a virtual screen. The rays
// produced are specified via a coordiate system on the virtual screen.  The
// left side of the virtual screen has x coordinate -1, the right side of the
// virtual screen has coordinate +1. The top of the virtual screen has y
// coordinate v and the bottom of the virtual screen has y coordinate -v (where
// the value of v depends on the aspect ratio of the screen).
type Camera interface {
	MakeRay(x, y float64) ray
}

type RectilinearCamera struct {
	screen, eye struct {
		// X vectors go from the center of the screen or eye to the right of
		// the screen or eye.  Y vectors go from the center of the screen or
		// eye towards the top of the screen or eye.  Loc is the location of
		// the center of the screen or eye.
		loc, x, y Vect
	}
}

// CameraConfig gives configuration options that are common to all camera
// types. This struct is a parameter to the camera factory functions, primarily
// so in the calling context, it's clear (from the name) whach each
// configuration option is (compared to just passing in Vects and float64s to
// the factory function).
type CameraConfig struct {
	Location      Vect
	ViewDirection Vect
	UpDirection   Vect
	FieldOfView   float64
	FocalLength   float64
	FocalRatio    float64
}

func NewRectilinearCamera(conf CameraConfig) Camera {

	cam := &RectilinearCamera{}

	conf.UpDirection = conf.UpDirection.Unit()
	conf.ViewDirection = conf.ViewDirection.Unit()

	cam.screen.x = Cross(conf.ViewDirection, conf.UpDirection)
	cam.screen.y = Cross(cam.screen.x, conf.ViewDirection)

	cam.eye.x = cam.screen.x.Extended(conf.FocalLength / conf.FocalRatio)
	cam.eye.y = cam.screen.y.Extended(conf.FocalLength / conf.FocalRatio)
	cam.eye.loc = conf.Location

	halfScreenWidth := math.Tan(conf.FieldOfView/2) * conf.FocalLength
	cam.screen.x = cam.screen.x.Extended(halfScreenWidth)
	cam.screen.y = cam.screen.y.Extended(halfScreenWidth)
	cam.screen.loc = Add(cam.eye.loc, conf.ViewDirection.Extended(conf.FocalLength))

	return cam
}

func (c *RectilinearCamera) MakeRay(x, y float64) ray {
	return ray{
	//start: c.eye.loc.Plus().Plus(...)
	}
}
