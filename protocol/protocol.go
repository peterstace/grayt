package protocol

import (
	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/xmath"
)

type Material struct {
	Colour    colour.Colour `json:"colour"`
	Emittance float64       `json:"emittance"`
	Mirror    bool          `json:"mirror"`
}

type Object struct {
	Surface  interface{} `json:"surface"`
	Material Material    `json:"material"`
}

type Camera struct {
	Location             xmath.Vector `json:"location"`
	LookingAt            xmath.Vector `json:"looking_at"`
	UpDirection          xmath.Vector `json:"up_direction"`
	FieldOfViewInRadians float64      `json:"field_of_view_in_radians"`
	FocalLength          float64      `json:"focal_length"`
	FocalRatio           float64      `json:"focal_ratio"`
	AspectWide           int          `json:"aspect_wide"`
	AspectHigh           int          `json:"aspect_high"`
}

type Scene struct {
	Camera  Camera
	Objects []Object
}

type Triangle struct {
	A xmath.Vector `json:"a"`
	B xmath.Vector `json:"b"`
	C xmath.Vector `json:"c"`
}

type AlignedBox struct {
	CornerA xmath.Vector `json:"a"`
	CornerB xmath.Vector `json:"b"`
}

type Sphere struct {
	Center xmath.Vector `json:"center"`
	Radius float64      `json:"radius"`
}

type AlignXSquare struct {
	X  float64 `json:"x"`
	Y1 float64 `json:"y_1"`
	Y2 float64 `json:"y_2"`
	Z1 float64 `json:"z_1"`
	Z2 float64 `json:"z_2"`
}

type AlignYSquare struct {
	X1 float64 `json:"x_1"`
	X2 float64 `json:"x_2"`
	Y  float64 `json:"y"`
	Z1 float64 `json:"z_1"`
	Z2 float64 `json:"z_2"`
}

type AlignZSquare struct {
	X1 float64 `json:"x_1"`
	X2 float64 `json:"x_2"`
	Y1 float64 `json:"y_1"`
	Y2 float64 `json:"y_2"`
	Z  float64 `json:"z"`
}

type Disc struct {
	Center xmath.Vector `json:"center"`
	Radius float64      `json:"radius"`
	// TODO: Should just accept any norm.
	UnitNorm xmath.Vector `json:"unit_norm"`
}

type Pipe struct {
	EndpointA xmath.Vector `json:"endpoint_a"`
	EndpointB xmath.Vector `json:"endpoint_b"`
	Radius    float64      `json:"radius"`
}
