package scene

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Camera struct {
	Location             Vector
	ViewDirection        Vector
	UpDirection          Vector
	FieldOfViewInDegrees float64
	FocalLength          float64 // Distance to the focus plane.
	FocalRatio           float64 // Ratio between the focal length and the diameter of the eye.
}

func (c Camera) String() string {
	return fmt.Sprintf("Location=%v ViewDir=%v UpDir=%v FOV=%v FocalLength=%v FocalRation=%v",
		c.Location, c.ViewDirection, c.UpDirection, c.FieldOfViewInDegrees, c.FocalLength, c.FocalRatio)
}

func DefaultCamera() Camera {
	return Camera{
		Location:             Vector{},
		ViewDirection:        Vect(0, 0, -1),
		UpDirection:          Vect(0, 1, 0),
		FieldOfViewInDegrees: 100,
		FocalLength:          10,
		FocalRatio:           math.Inf(+1),
	}
}

type Colour [3]float64

var White = Colour{1, 1, 1}

type Triangle struct {
	A, B, C   Vector
	Colour    Colour
	Emittance float64
}

func (t Triangle) String() string {
	return fmt.Sprintf("A=%v B=%v C=%v Colour=%v Emittance=%v",
		t.A, t.B, t.C, t.Colour, t.Emittance)
}

func JoinTriangles(ts ...[]Triangle) []Triangle {
	var all []Triangle
	for _, t := range ts {
		all = append(all, t...)
	}
	return all
}

type Vector struct {
	X, Y, Z float64
}

func Vect(x, y, z float64) Vector {
	return Vector{x, y, z}
}

func (v Vector) Midpoint(u Vector) Vector {
	return v.Add(u).Scale(0.5)
}

func (v Vector) Add(u Vector) Vector {
	return Vector{v.X + u.X, v.Y + u.Y, v.Z + u.Z}
}

func (v Vector) Scale(f float64) Vector {
	return Vector{v.X * f, v.Y * f, v.Z * f}
}

type Scene struct {
	Camera    Camera
	Triangles []Triangle
}

func (s Scene) WriteTo(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, s.Camera); err != nil {
		return err
	}
	for _, tri := range s.Triangles {
		if err := binary.Write(w, binary.BigEndian, tri); err != nil {
			return err
		}
	}
	return nil
}

func ReadFrom(r io.Reader) (Scene, error) {
	var s Scene
	if err := binary.Read(r, binary.BigEndian, &s.Camera); err != nil {
		return Scene{}, err
	}
	var tri Triangle
	for {
		switch err := binary.Read(r, binary.BigEndian, &tri); err {
		case io.EOF:
			return s, nil
		case nil:
			s.Triangles = append(s.Triangles, tri)
		default:
			return Scene{}, err
		}
	}
}
