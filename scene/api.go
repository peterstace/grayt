package scene

import "math"

/*
	The following types describe a scene.
*/

type Triangle struct {
	A, B, C   Vector
	Colour    Colour
	Emittance float64
}

type Camera struct {
	Location      Vector
	ViewDirection Vector
	UpDirection   Vector
	FieldOfView   float64 // In degrees.
	FocalLength   float64 // Distance to the focus plane.
	FocalRatio    float64 // Ratio between the focal length and the diameter of the eye.
}

func DefaultCamera() CameraDescription {
	return CameraDescription{
		Location:      Vector{},
		ViewDirection: Vect(0, 0, -1),
		UpDirection:   Vect(0, 1, 0),
		FieldOfView:   100,
		FocalLength:   10,
		FocalRatio:    math.Inf(+1),
	}
}
