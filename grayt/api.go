package grayt

import (
	"fmt"
	"math"
)

type CameraBlueprint struct {
	location             Vector
	lookingAt            Vector
	upDirection          Vector
	fieldOfViewInRadians float64
	focalLength          float64
	focalRatio           float64
	aspectWide           int
	aspectHigh           int
}

func (c CameraBlueprint) With(opts ...cameraOption) CameraBlueprint {
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

func (c CameraBlueprint) String() string {
	return fmt.Sprintf(
		"Location=%v LookingAt=%v UpDir=%v FOV=%v FocalLength=%v FocalRatio=%v",
		c.location, c.lookingAt, c.upDirection, c.fieldOfViewInRadians, c.focalLength, c.focalRatio,
	)
}

func Camera() CameraBlueprint {
	return CameraBlueprint{
		location:             Vect(0, 10, 10),
		lookingAt:            Vect(0, 0, 0),
		upDirection:          Vect(0, 1, 0),
		fieldOfViewInRadians: 90 * math.Pi / 180,
		focalLength:          1.0,
		focalRatio:           math.Inf(+1),
		aspectWide:           1,
		aspectHigh:           1,
	}
}

type cameraOption func(*CameraBlueprint)

func Location(x Vector) cameraOption {
	return func(c *CameraBlueprint) {
		c.location = x
	}
}

func LookingAt(x Vector) cameraOption {
	return func(c *CameraBlueprint) {
		c.lookingAt = x
	}
}

func UpDirection(x Vector) cameraOption {
	return func(c *CameraBlueprint) {
		c.upDirection = x
	}
}

func FieldOfViewInRadians(r float64) cameraOption {
	return func(c *CameraBlueprint) {
		c.fieldOfViewInRadians = r
	}
}

func FieldOfViewInDegrees(d float64) cameraOption {
	return FieldOfViewInRadians(d * math.Pi / 180)
}

func ScaleFieldOfView(s float64) cameraOption {
	return func(c *CameraBlueprint) {
		c.fieldOfViewInRadians *= s
	}
}

func FocalLengthAndRatio(focalLength, focalRatio float64) cameraOption {
	return func(c *CameraBlueprint) {
		c.focalLength = focalLength
		c.focalRatio = focalRatio
	}
}

func AspectRatioWidthAndHeight(wide, high int) cameraOption {
	if wide <= 0 || high <= 0 {
		panic("aspect ratio elements must be positive")
	}
	return func(c *CameraBlueprint) {
		c.aspectWide = wide
		c.aspectHigh = high
	}
}

type Scene struct {
	Camera  CameraBlueprint
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
	White   = 0xffffff
	Black   = 0x000000
	Red     = 0xff0000
	Green   = 0x00ff00
	Blue    = 0x0000ff
	Yellow  = 0xffff00
	Cyan    = 0x00ffff
	Magenta = 0xff00ff
)

func ColourRGB(rgb uint32) func(*Object) {
	return func(o *Object) {
		o.material.colour = newColourFromRGB(rgb)
	}
}

func ColourHSL(hue, saturation, lightness float64) func(*Object) {
	return func(o *Object) {
		o.material.colour = newColourFromHSL(hue, saturation, lightness)
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

	same := func(a, b float64) int {
		if a == b {
			return 1
		}
		return 0
	}
	if same(a.X, b.X)+same(a.Y, b.Y)+same(a.Z, b.Z) != 1 {
		panic("a and b must have exactly 1 dimension in common")
	}

	a, b = a.Min(b), a.Max(b)

	switch {
	case a.X == b.X:
		return defaultObject(&alignXSquare{a.X, a.Y, b.Y, a.Z, b.Z})
	case a.Y == b.Y:
		return defaultObject(&alignYSquare{a.X, b.X, a.Y, a.Z, b.Z})
	case a.Z == b.Z:
		return defaultObject(&alignZSquare{a.X, b.X, a.Y, b.Y, a.Z})
	default:
		panic(false)

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
	return defaultObject(&plane{normal, pt})
}

func XPlane(x float64) ObjectList {
	return defaultObject(&alignXPlane{x})
}

func YPlane(y float64) ObjectList {
	return defaultObject(&alignYPlane{y})
}

func ZPlane(z float64) ObjectList {
	return defaultObject(&alignZPlane{z})
}

func Sphere(c Vector, r float64) ObjectList {
	return defaultObject(&sphere{c, r})
}

func defaultObject(s surface) ObjectList {
	return ObjectList{{
		s, material{colour: newColourFromRGB(White)},
	}}
}
