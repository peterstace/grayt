package trace

/*
type CameraBlueprint struct {
	Location             xmath.Vector `json:"location"`
	LookingAt            xmath.Vector `json:"looking_at"`
	UpDirection          xmath.Vector `json:"up_direction"`
	FieldOfViewInRadians float64      `json:"field_of_view_in_radians"`
	FocalLength          float64      `json:"focal_length"`
	FocalRatio           float64      `json:"focal_ratio"`
	AspectWide           int          `json:"aspect_wide"`
	AspectHigh           int          `json:"aspect_high"`
}
*/

/*
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

func Camera() CameraBlueprint {
	return CameraBlueprint{
		Location:             xmath.Vect(0, 10, 10),
		LookingAt:            xmath.Vect(0, 0, 0),
		UpDirection:          xmath.Vect(0, 1, 0),
		FieldOfViewInRadians: 90 * math.Pi / 180,
		FocalLength:          1.0,
		FocalRatio:           math.MaxFloat64,
		AspectWide:           1,
		AspectHigh:           1,
	}
}
*/

/*
type cameraOption func(*CameraBlueprint)

func Location(x xmath.Vector) cameraOption {
	return func(c *CameraBlueprint) {
		c.Location = x
	}
}

func LookingAt(x xmath.Vector) cameraOption {
	return func(c *CameraBlueprint) {
		c.LookingAt = x
	}
}

func UpDirection(x xmath.Vector) cameraOption {
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
*/

//type ObjectList []Object

/*
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
*/

/*
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
*/
/*
func ColourRGB(rgb uint32) func(*Object) {
	return func(o *Object) {
		o.Material.Colour = colour.NewColourFromRGB(rgb)
	}
}

func ColourHSL(hue, saturation, lightness float64) func(*Object) {
	return func(o *Object) {
		o.Material.Colour = colour.NewColourFromHSL(hue, saturation, lightness)
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

func Translate(v xmath.Vector) func(*Object) {
	return func(o *Object) {
		o.Surface.translate(v)
	}
}

func RotateDegrees(v xmath.Vector, degs float64) func(*Object) {
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

func Triangle(a, b, c xmath.Vector) ObjectList {
	return defaultObject(newTriangle(a, b, c))
}

func AlignedBox(a, b xmath.Vector) ObjectList {
	return defaultObject(newAlignedBox(a, b))
}

func Sphere(c xmath.Vector, r float64) ObjectList {
	return defaultObject(&sphere{c, r})
}

func Disc(c xmath.Vector, r float64, n xmath.Vector) ObjectList {
	return defaultObject(&disc{c, r * r, n.Unit()})
}

func Pipe(a, b xmath.Vector, r float64) ObjectList {
	return defaultObject(&pipe{a, b, r})
}

func defaultObject(s surface) ObjectList {
	return ObjectList{{
		s, material{Colour: colour.NewColourFromRGB(White)},
	}}
}
*/
