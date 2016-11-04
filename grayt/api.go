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

type Scene struct {
	Camera  Camera
	Objects ObjectList
}

type ObjectList []Object

func Group(objLists ...ObjectList) ObjectList {
	var grouped ObjectList
	for _, objList := range objLists {
		grouped = append(grouped, objList...)
	}
	return grouped
}

func (o ObjectList) With(fns ...func(*Object)) ObjectList {
	for i := range o {
		for _, fn := range fns {
			fn(&o[i])
		}
	}
	return o
}

const (
	White = 0xffffff
	Black = 0x000000
	Red   = 0xff0000
	Green = 0x00ff00
	Blue  = 0x0000ff
)

func ColourRGB(rgb uint32) func(*Object) {
	return func(o *Object) {
		o.material.colour = newColour(rgb)
	}
}

func Emittance(e float64) func(*Object) {
	return func(o *Object) {
		o.material.emittance = e
	}
}

func Triangle(a, b, c Vector) ObjectList {
	return ObjectList{{
		newTriangle(a, b, c),
		material{colour: newColour(White)},
	}}
}

func AlignedSquare(a, b Vector) ObjectList {
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
	return Group(
		Triangle(a, c, d),
		Triangle(b, c, d),
	)
}

func AlignedBox(a, b Vector) ObjectList {

	a1 := Vect(b.X, a.Y, a.Z)
	a2 := Vect(a.X, b.Y, a.Z)
	a3 := Vect(a.X, a.Y, b.Z)
	b1 := Vect(a.X, b.Y, b.Z)
	b2 := Vect(b.X, a.Y, b.Z)
	b3 := Vect(b.X, b.Y, a.Z)

	return Group(
		Triangle(a, a1, a2),
		Triangle(a, a2, a3),
		Triangle(a, a3, a1),

		Triangle(b, b1, b2),
		Triangle(b, b2, b3),
		Triangle(b, b3, b1),

		Triangle(a1, b2, b3),
		Triangle(a2, b3, b1),
		Triangle(a3, b1, b2),

		Triangle(b1, a2, a3),
		Triangle(b2, a3, a1),
		Triangle(b3, a1, a2),
	)
}

func Square(a, b, c, d Vector) ObjectList {
	return Group(
		Triangle(a, b, c),
		Triangle(c, d, a),
	)
}
