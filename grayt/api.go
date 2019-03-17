package grayt

import (
	"fmt"
	"math"
)

type CameraBlueprint struct {
	Location             Vector  `json:"location"`
	LookingAt            Vector  `json:"looking_at"`
	UpDirection          Vector  `json:"up_direction"`
	FieldOfViewInRadians float64 `json:"field_of_view_in_radians"`
	FocalLength          float64 `json:"focal_length"`
	FocalRatio           float64 `json:"focal_ratio"`
	AspectWide           int     `json:"aspect_wide"`
	AspectHigh           int     `json:"aspect_high"`
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
		c.Location, c.LookingAt, c.UpDirection, c.FieldOfViewInRadians, c.FocalLength, c.FocalRatio,
	)
}

// TODO: Can get rid of this?
func (c CameraBlueprint) pxHigh(pxWide int) int {
	return pxWide * c.AspectHigh / c.AspectWide
}

func Camera() CameraBlueprint {
	return CameraBlueprint{
		Location:             Vect(0, 10, 10),
		LookingAt:            Vect(0, 0, 0),
		UpDirection:          Vect(0, 1, 0),
		FieldOfViewInRadians: 90 * math.Pi / 180,
		FocalLength:          1.0,
		FocalRatio:           math.MaxFloat64,
		AspectWide:           1,
		AspectHigh:           1,
	}
}

type cameraOption func(*CameraBlueprint)

func Location(x Vector) cameraOption {
	return func(c *CameraBlueprint) {
		c.Location = x
	}
}

func LookingAt(x Vector) cameraOption {
	return func(c *CameraBlueprint) {
		c.LookingAt = x
	}
}

func UpDirection(x Vector) cameraOption {
	return func(c *CameraBlueprint) {
		c.UpDirection = x
	}
}

func FieldOfViewInRadians(r float64) cameraOption {
	return func(c *CameraBlueprint) {
		c.FieldOfViewInRadians = r
	}
}

func FieldOfViewInDegrees(d float64) cameraOption {
	return FieldOfViewInRadians(d * math.Pi / 180)
}

func ScaleFieldOfView(s float64) cameraOption {
	return func(c *CameraBlueprint) {
		c.FieldOfViewInRadians *= s
	}
}

func FocalLengthAndRatio(focalLength, focalRatio float64) cameraOption {
	return func(c *CameraBlueprint) {
		c.FocalLength = focalLength
		c.FocalRatio = focalRatio
	}
}

func AspectRatioWidthAndHeight(wide, high int) cameraOption {
	if wide <= 0 || high <= 0 {
		panic("aspect ratio elements must be positive")
	}
	return func(c *CameraBlueprint) {
		c.AspectWide = wide
		c.AspectHigh = high
	}
}

type Scene struct {
	Camera  CameraBlueprint `json:"camera"`
	Objects ObjectList      `json:"objects"`
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
		o.Material.Colour = newColourFromRGB(rgb)
	}
}

func ColourHSL(hue, saturation, lightness float64) func(*Object) {
	return func(o *Object) {
		o.Material.Colour = newColourFromHSL(hue, saturation, lightness)
	}
}

func Emittance(e float64) func(*Object) {
	return func(o *Object) {
		o.Material.Emittance = e
	}
}

func Mirror() func(*Object) {
	return func(o *Object) {
		o.Material.Mirror = true
	}
}

func Translate(v Vector) func(*Object) {
	return func(o *Object) {
		o.Surface.translate(v)
	}
}

func RotateDegrees(v Vector, degs float64) func(*Object) {
	return func(o *Object) {
		o.Surface.rotate(v, degs/180*math.Pi)
	}
}

func Scale(f float64) func(*Object) {
	return func(o *Object) {
		o.Surface.scale(f)
	}
}

func BoundingBox() func(*Object) {
	return func(o *Object) {
		a, b := o.Surface.bound()
		o.Surface = newAlignedBox(a, b)
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

func Sphere(c Vector, r float64) ObjectList {
	return defaultObject(&sphere{c, r})
}

func Disc(c Vector, r float64, n Vector) ObjectList {
	return defaultObject(&disc{c, r * r, n.Unit()})
}

func Pipe(a, b Vector, r float64) ObjectList {
	return defaultObject(&pipe{a, b, r})
}

func defaultObject(s surface) ObjectList {
	return ObjectList{{
		s, material{Colour: newColourFromRGB(White)},
	}}
}
