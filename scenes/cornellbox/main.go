package main

import (
	"math"

	"github.com/peterstace/grayt"
)

func main() {
	grayt.NewRunner().Run(scene())
}

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
	up    = grayt.Vect(0.0, 1.0, 0.0)
	down  = grayt.Vect(0.0, -1.0, 0.0)
	left  = grayt.Vect(-1.0, 0.0, 0.0)
	right = grayt.Vect(1.0, 0.0, 0.0)
	back  = grayt.Vect(0.0, 0.0, 1.0)
	zero  = grayt.Vect(0.0, 0.0, 0.0)
	one   = grayt.Vect(1.0, 1.0, -1.0)
)

var (
	white = grayt.Material{Color: grayt.Colour{1, 1, 1}, Emmitance: 0.0}
	green = grayt.Material{Color: grayt.Colour{0, 1, 0}, Emmitance: 0.0}
	red   = grayt.Material{Color: grayt.Colour{1, 0, 0}, Emmitance: 0.0}
	blue  = grayt.Material{Color: grayt.Colour{0, 0, 1}, Emmitance: 0.0}
)

func scene() grayt.Scene {
	ee := []grayt.Entity{{
		grayt.Material{grayt.Colour{1, 1, 1}, 5},
		[]grayt.SurfaceFactory{grayt.Sphere{grayt.Vect(0.5, 1.0, -0.5), 0.25}},
	}}
	ee = append(ee, box()...)
	ee = append(ee, grayt.Entity{white, tallBlock()})
	ee = append(ee, grayt.Entity{white, shortBlock()})
	return grayt.Scene{cam(), ee}
}

func cam() grayt.CameraConfig {
	const D = 1.3 // Estimated.
	return grayt.CameraConfig{
		Projection:    grayt.Rectilinear,
		Location:      grayt.Vect(0.5, 0.5, D),
		ViewDirection: grayt.Vect(0.0, 0.0, -1.0),
		UpDirection:   up,
		FieldOfView:   2 * math.Asin(0.5/math.Sqrt(0.25+D*D)),
		FocalLength:   0.5 + D,
		FocalRatio:    math.MaxFloat64, //math.Inf(1),
	}
}

func box() []grayt.Entity {
	return []grayt.Entity{
		{white, []grayt.SurfaceFactory{grayt.Plane{up, zero}}},
		{white, []grayt.SurfaceFactory{grayt.Plane{down, one}}},
		{red, []grayt.SurfaceFactory{grayt.Plane{right, zero}}},
		{green, []grayt.SurfaceFactory{grayt.Plane{left, one}}},
		{white, []grayt.SurfaceFactory{grayt.Plane{back, one}}},
		{white, []grayt.SurfaceFactory{grayt.AlignedBox{grayt.Vect(0.1, 0.0, -0.1), grayt.Vect(0.9, 0.1, -0.9)}}},
	}
}

func shortBlock() []grayt.SurfaceFactory {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = grayt.Vect(0.76, 0.00, -0.12)
		LBB = grayt.Vect(0.85, 0.00, -0.41)
		RBF = grayt.Vect(0.47, 0.00, -0.21)
		RBB = grayt.Vect(0.56, 0.00, -0.49)
		LTF = grayt.Vect(0.76, 0.30, -0.12)
		LTB = grayt.Vect(0.85, 0.30, -0.41)
		RTF = grayt.Vect(0.47, 0.30, -0.21)
		RTB = grayt.Vect(0.56, 0.30, -0.49)
	)
	return []grayt.SurfaceFactory{
		grayt.Square{LTF, LTB, RTB, RTF},
		grayt.Square{LBF, RBF, RTF, LTF},
		grayt.Square{LBB, RBB, RTB, LTB},
		grayt.Square{LBF, LBB, LTB, LTF},
		grayt.Square{RBF, RBB, RTB, RTF},
	}
}

func tallBlock() []grayt.SurfaceFactory {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = grayt.Vect(0.52, 0.00, -0.54)
		LBB = grayt.Vect(0.43, 0.00, -0.83)
		RBF = grayt.Vect(0.23, 0.00, -0.45)
		RBB = grayt.Vect(0.14, 0.00, -0.74)
		LTF = grayt.Vect(0.52, 0.60, -0.54)
		LTB = grayt.Vect(0.43, 0.60, -0.83)
		RTF = grayt.Vect(0.23, 0.60, -0.45)
		RTB = grayt.Vect(0.14, 0.60, -0.74)
	)
	return []grayt.SurfaceFactory{
		grayt.Square{LTF, LTB, RTB, RTF},
		grayt.Square{LBF, RBF, RTF, LTF},
		grayt.Square{LBB, RBB, RTB, LTB},
		grayt.Square{LBF, LBB, LTB, LTF},
		grayt.Square{RBF, RBB, RTB, RTF},
	}
}
