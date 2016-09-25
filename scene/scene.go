package scene

import (
	"encoding/binary"
	"io"
	"math"
)

type Camera struct {
	Location      Vector
	ViewDirection Vector
	UpDirection   Vector
	FieldOfView   float64 // In degrees.
	FocalLength   float64 // Distance to the focus plane.
	FocalRatio    float64 // Ratio between the focal length and the diameter of the eye.
}

func DefaultCamera() Camera {
	return Camera{
		Location:      Vector{},
		ViewDirection: Vect(0, 0, -1),
		UpDirection:   Vect(0, 1, 0),
		FieldOfView:   100,
		FocalLength:   10,
		FocalRatio:    math.Inf(+1),
	}
}

type Colour [3]uint16

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

func (s Scene) WriteTo(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, s)
}

func ReadFrom(r io.Reader) (Scene, error) {
	var s Scene
	err := binary.Read(r, binary.BigEndian, &s)
	return s, err
}
