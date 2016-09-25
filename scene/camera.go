package scene

import (
	"io"
	"math"
)

type CameraDescription struct {
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

type Triangle struct {
	A, B, C   Vector
	Colour    Colour
	Emittance float64
}
type Vector struct {
	X, Y, Z float64
}

func Vect(x, y, z float64) Vector {
	return Vector{x, y, z}
}

type Scene struct {
	Camera    Camera
	Triangles []Triangle
}

func (s Scene) WriteTo(w io.Writer) (n int64, err error) {
	// TODO
	return 0, nil
}

func ReadFrom(r io.Reader) (Scene, error) {
	// TODO
	return Scene{}, nil
}
