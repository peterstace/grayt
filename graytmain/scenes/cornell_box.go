package scenes

import (
	"math"

	"github.com/peterstace/grayt/graytlib"
)

// Cornell Box
//
// (0,y,-1) +---------+ (1,y,-1)
//          | BB      |
//          | BB  BB  |
//          |     BB  |
// (0,y,0)  +         + (1,y,0)
//            \     /
//             \   /
//              \ /
//               C (0.5,0.5,D)
//
//
// SINE(0.5*T) = 0.5 / SQRT(0.5^2 + D^2)
// T = 2*ARCSINE(0.5/SQRT(0.25 + D^2))

type CornellBox struct{}

func (c CornellBox) Name() string { return "Cornell Box" }

var (
	up    = graytlib.Vect{0.0, 1.0, 0.0}
	down  = graytlib.Vect{0.0, -1.0, 0.0}
	left  = graytlib.Vect{-1.0, 0.0, 0.0}
	right = graytlib.Vect{1.0, 0.0, 0.0}
	back  = graytlib.Vect{0.0, 0.0, 1.0}
	zero  = graytlib.Vect{0.0, 0.0, 0.0}
	one   = graytlib.Vect{1.0, 1.0, -1.0}
)

var (
	white = graytlib.Material{graytlib.Colour{1, 1, 1}, 0.0}
	green = graytlib.Material{graytlib.Colour{0, 1, 0}, 0.0}
	red   = graytlib.Material{graytlib.Colour{1, 0, 0}, 0.0}
	blue  = graytlib.Material{graytlib.Colour{0, 0, 1}, 0.0}
)

func cam() graytlib.CameraConfig {
	const D = 1.3 // Estimated.
	return graytlib.CameraConfig{
		Projection:    graytlib.Rectilinear,
		Location:      graytlib.Vect{0.5, 0.5, D},
		ViewDirection: graytlib.Vect{0.0, 0.0, -1.0},
		UpDirection:   up,
		FieldOfView:   2 * math.Asin(0.5/math.Sqrt(0.25+D*D)),
		FocalLength:   0.5 + D,
		FocalRatio:    math.MaxFloat64, //math.Inf(1),
	}
}

func (c CornellBox) Scene() graytlib.Scene {
	ee := []graytlib.Entity{{
		graytlib.Material{graytlib.Colour{1, 1, 1}, 5},
		[]graytlib.SurfaceFactory{graytlib.Sphere{graytlib.Vect{0.5, 1.0, -0.5}, 0.25}},
	}}
	ee = append(ee, box()...)
	ee = append(ee, graytlib.Entity{white, tallBlock()})
	ee = append(ee, graytlib.Entity{white, shortBlock()})
	return graytlib.Scene{cam(), ee}
}

func box() []graytlib.Entity {
	return []graytlib.Entity{
		{white, []graytlib.SurfaceFactory{graytlib.Plane{up, zero}}},
		{white, []graytlib.SurfaceFactory{graytlib.Plane{down, one}}},
		{red, []graytlib.SurfaceFactory{graytlib.Plane{right, zero}}},
		{green, []graytlib.SurfaceFactory{graytlib.Plane{left, one}}},
		{white, []graytlib.SurfaceFactory{graytlib.Plane{back, one}}},
		{white, []graytlib.SurfaceFactory{graytlib.AlignedBox{graytlib.Vect{0.1, 0.0, -0.1}, graytlib.Vect{0.9, 0.1, -0.9}}}},
	}
}

func shortBlock() []graytlib.SurfaceFactory {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = graytlib.Vect{0.76, 0.00, -0.12}
		LBB = graytlib.Vect{0.85, 0.00, -0.41}
		RBF = graytlib.Vect{0.47, 0.00, -0.21}
		RBB = graytlib.Vect{0.56, 0.00, -0.49}
		LTF = graytlib.Vect{0.76, 0.30, -0.12}
		LTB = graytlib.Vect{0.85, 0.30, -0.41}
		RTF = graytlib.Vect{0.47, 0.30, -0.21}
		RTB = graytlib.Vect{0.56, 0.30, -0.49}
	)
	return []graytlib.SurfaceFactory{
		graytlib.Square{LTF, LTB, RTB, RTF},
		graytlib.Square{LBF, RBF, RTF, LTF},
		graytlib.Square{LBB, RBB, RTB, LTB},
		graytlib.Square{LBF, LBB, LTB, LTF},
		graytlib.Square{RBF, RBB, RTB, RTF},
	}
}

func tallBlock() []graytlib.SurfaceFactory {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = graytlib.Vect{0.52, 0.00, -0.54}
		LBB = graytlib.Vect{0.43, 0.00, -0.83}
		RBF = graytlib.Vect{0.23, 0.00, -0.45}
		RBB = graytlib.Vect{0.14, 0.00, -0.74}
		LTF = graytlib.Vect{0.52, 0.60, -0.54}
		LTB = graytlib.Vect{0.43, 0.60, -0.83}
		RTF = graytlib.Vect{0.23, 0.60, -0.45}
		RTB = graytlib.Vect{0.14, 0.60, -0.74}
	)
	return []graytlib.SurfaceFactory{
		graytlib.Square{LTF, LTB, RTB, RTF},
		graytlib.Square{LBF, RBF, RTF, LTF},
		graytlib.Square{LBB, RBB, RTB, LTB},
		graytlib.Square{LBF, LBB, LTB, LTF},
		graytlib.Square{RBF, RBB, RTB, RTF},
	}
}
