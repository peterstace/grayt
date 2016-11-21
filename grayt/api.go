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
	return defaultObject(newTriangle(a, b, c))
}

func AlignedSquare(a, b Vector) ObjectList {
	a, b = a.min(b), a.max(b)
	switch {
	case a.X == b.X:
		return defaultObject(&alignXSquare{a.X, a.Y, b.Y, a.Z, b.Z})
	case a.Y == b.Y:
		return defaultObject(&alignYSquare{a.X, b.X, a.Y, a.Z, b.Z})
	case a.Z == b.Z:
		return defaultObject(&alignZSquare{a.X, b.X, a.Y, b.Y, a.Z})
	default:
		panic("a and b lie in a common aligned plane")

	}
}

func AlignedBox(a, b Vector) ObjectList {
	return defaultObject(newAlignedBox(a, b))
}

func Square(a, b, c, d Vector) ObjectList {
	return Group(
		Triangle(a, b, c),
		Triangle(c, d, a),
	)
}

func Plane(normal, pt Vector) ObjectList {
	return defaultObject(newPlane(normal, pt))
}

func XPlane(x float64) ObjectList {
	return defaultObject(newXPlane(x))
}

func YPlane(y float64) ObjectList {
	return defaultObject(newYPlane(y))
}

func ZPlane(z float64) ObjectList {
	return defaultObject(newZPlane(z))
}

func Sphere(c Vector, r float64) ObjectList {
	return defaultObject(newSphere(c, r))
}

func defaultObject(s surface) ObjectList {
	return ObjectList{{
		s, material{colour: newColour(White)},
	}}
}
