package grayt

import (
	"fmt"
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

type Triangle struct {
	A, B, C   Vector
	Colour    Colour
	Emittance float64
}

func (t Triangle) String() string {
	return fmt.Sprintf("A=%v B=%v C=%v Colour=%v Emittance=%v",
		t.A, t.B, t.C, t.Colour, t.Emittance)
}

type TriangleList []Triangle

func JoinTriangles(ts ...TriangleList) TriangleList {
	var all []Triangle
	for _, t := range ts {
		all = append(all, t...)
	}
	return all
}

func (t TriangleList) SetColour(c Colour) TriangleList {
	for i := range t {
		t[i].Colour = c
	}
	return t
}

func (t TriangleList) SetEmittance(e float64) TriangleList {
	for i := range t {
		t[i].Emittance = e
	}
	return t
}

func Square(a, b, c, d Vector) TriangleList {
	return []Triangle{
		{A: a, B: b, C: c},
		{A: c, B: d, C: a},
	}
}

func AlignedSquare(a, b Vector) TriangleList {
	var c, d Vector
	switch {
	case a.X == b.X:
		c = Vector{X: a.X, Y: a.Y, Z: b.Z}
		d = Vector{X: a.X, Y: b.Y, Z: a.Z}
	case a.Y == b.Y:
		c = Vector{X: a.X, Y: a.Y, Z: b.Z}
		d = Vector{X: b.X, Y: a.Y, Z: a.Z}
	case a.Z == b.Z:
		c = Vector{X: a.X, Y: b.Y, Z: a.Z}
		d = Vector{X: b.X, Y: a.Y, Z: a.Z}
	default:
		panic("a and b line in a common aligned plane")

	}
	return []Triangle{
		{A: a, B: c, C: d},
		{A: b, B: c, C: d},
	}
}

func AlignedBox(a, b Vector) TriangleList {
	a1 := Vect(b.X, a.Y, a.Z)
	a2 := Vect(a.X, b.Y, a.Z)
	a3 := Vect(a.X, a.Y, b.Z)
	b1 := Vect(a.X, b.Y, b.Z)
	b2 := Vect(b.X, a.Y, b.Z)
	b3 := Vect(b.X, b.Y, a.Z)

	return TriangleList{

		{A: a, B: a1, C: a2},
		{A: a, B: a2, C: a3},
		{A: a, B: a3, C: a1},

		{A: b, B: b1, C: b2},
		{A: b, B: b2, C: b3},
		{A: b, B: b3, C: b1},

		{A: a1, B: b2, C: b3},
		{A: a2, B: b3, C: b1},
		{A: a3, B: b1, C: b2},

		{A: b1, B: a2, C: a3},
		{A: b2, B: a3, C: a1},
		{A: b3, B: a1, C: a2},
	}
}

type Scene struct {
	Camera    Camera
	Triangles []Triangle
}
