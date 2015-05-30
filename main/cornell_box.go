package main

import (
	"math"

	"github.com/peterstace/grayt"
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

var (
	up    = grayt.Vect{0.0, 1.0, 0.0}
	down  = grayt.Vect{0.0, -1.0, 0.0}
	left  = grayt.Vect{-1.0, 0.0, 0.0}
	right = grayt.Vect{1.0, 0.0, 0.0}
	back  = grayt.Vect{0.0, 0.0, 1.0}
	zero  = grayt.Vect{0.0, 0.0, 0.0}
	one   = grayt.Vect{1.0, 1.0, -1.0}
)

var (
	white = grayt.Material{grayt.Colour{1, 1, 1}}
	green = grayt.Material{grayt.Colour{0, 1, 0}}
	red   = grayt.Material{grayt.Colour{1, 0, 0}}
	blue  = grayt.Material{grayt.Colour{0, 0, 1}}
)

func cam() grayt.Camera {
	const D = 1.3
	return grayt.NewRectilinearCamera(grayt.CameraConfig{
		Location:      grayt.Vect{0.5, 0.5, D},
		ViewDirection: grayt.Vect{0.0, 0.0, -1.0},
		UpDirection:   up,
		FieldOfView:   2 * math.Asin(0.5/math.Sqrt(0.25+D*D)),
		FocalLength:   0.5 + D,
		FocalRatio:    10.0, //math.Inf(+1),
	})
}

func box() []grayt.Geometry {
	return []grayt.Geometry{
		grayt.NewPlane(white, up, zero),
		grayt.NewPlane(white, down, one),
		grayt.NewPlane(red, right, zero),
		grayt.NewPlane(green, left, one),
		grayt.NewPlane(white, back, one),
	}
}

func shortBlock(m grayt.Material) []grayt.Geometry {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = grayt.Vect{0.76, 0.00, -0.12}
		LBB = grayt.Vect{0.85, 0.00, -0.41}
		RBF = grayt.Vect{0.47, 0.00, -0.21}
		RBB = grayt.Vect{0.56, 0.00, -0.49}
		LTF = grayt.Vect{0.76, 0.30, -0.12}
		LTB = grayt.Vect{0.85, 0.30, -0.41}
		RTF = grayt.Vect{0.47, 0.30, -0.21}
		RTB = grayt.Vect{0.56, 0.30, -0.49}
	)
	var gs []grayt.Geometry
	gs = append(gs, grayt.NewSquare(m, LTF, LTB, RTB, RTF)...)
	gs = append(gs, grayt.NewSquare(m, LBF, RBF, RTF, LTF)...)
	gs = append(gs, grayt.NewSquare(m, LBB, RBB, RTB, LTB)...)
	gs = append(gs, grayt.NewSquare(m, LBF, LBB, LTB, LTF)...)
	gs = append(gs, grayt.NewSquare(m, RBF, RBB, RTB, RTF)...)
	return gs
}

func tallBlock(m grayt.Material) []grayt.Geometry {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = grayt.Vect{0.52, 0.00, -0.54}
		LBB = grayt.Vect{0.43, 0.00, -0.83}
		RBF = grayt.Vect{0.23, 0.00, -0.45}
		RBB = grayt.Vect{0.14, 0.00, -0.74}
		LTF = grayt.Vect{0.52, 0.60, -0.54}
		LTB = grayt.Vect{0.43, 0.60, -0.83}
		RTF = grayt.Vect{0.23, 0.60, -0.45}
		RTB = grayt.Vect{0.14, 0.60, -0.74}
	)
	var gs []grayt.Geometry
	gs = append(gs, grayt.NewSquare(m, LTF, LTB, RTB, RTF)...)
	gs = append(gs, grayt.NewSquare(m, LBF, RBF, RTF, LTF)...)
	gs = append(gs, grayt.NewSquare(m, LBB, RBB, RTB, LTB)...)
	gs = append(gs, grayt.NewSquare(m, LBF, LBB, LTB, LTF)...)
	gs = append(gs, grayt.NewSquare(m, RBF, RBB, RTB, RTF)...)
	return gs
}

func CornellBoxStandard() grayt.Scene {
	s := grayt.Scene{
		Camera: cam(),
		Lights: []grayt.Light{
			grayt.Light{
				Location:  grayt.Vect{0.5, 0.9, -0.5},
				Intensity: 0.3,
			},
		},
	}
	s.Geometries = append(s.Geometries, box()...)
	s.Geometries = append(s.Geometries, tallBlock(white)...)
	s.Geometries = append(s.Geometries, shortBlock(white)...)
	return s
}
