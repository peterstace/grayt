package tracer

import (
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
