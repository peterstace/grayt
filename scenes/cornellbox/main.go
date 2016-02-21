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
	white = grayt.Material{Color: grayt.White, Emmitance: 0.0}
	green = grayt.Material{Color: grayt.Green, Emmitance: 0.0}
	red   = grayt.Material{Color: grayt.Red, Emmitance: 0.0}
	blue  = grayt.Material{Color: grayt.Blue, Emmitance: 0.0}
)

func scene() grayt.Scene {
	ee := []grayt.Entity{{
		grayt.Material{Colour: grayt.White, Emittance: 5},
		[]grayt.SurfaceFactory{grayt.Sphere{C: grayt.Vect(0.5, 1.0, -0.5), R: 0.25}},
	}}
	ee = append(ee, box()...)
	ee = append(ee, grayt.Entity{Material: white, SurfaceFactories: tallBlock()})
	ee = append(ee, grayt.Entity{Material: white, SurfaceFactorie: shortBlock()})
	return grayt.Scene{CameraConfig: cam(), Entities: ee}
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
		{Material: white, Entities: []grayt.SurfaceFactory{grayt.Plane{N: up, X: zero}}},
		{Material: white, Entities: []grayt.SurfaceFactory{grayt.Plane{N: down, X: one}}},
		{Material: red, Entities: []grayt.SurfaceFactory{grayt.Plane{N: right, X: zero}}},
		{Material: green, Entities: []grayt.SurfaceFactory{grayt.Plane{N: left, X: one}}},
		{Material: white, Entities: []grayt.SurfaceFactory{grayt.Plane{N: back, X: one}}},
		{Material: white, Entities: []grayt.SurfaceFactory{grayt.AlignedBox{Corner1: grayt.Vect(0.1, 0.0, -0.1), Corner2: grayt.Vect(0.9, 0.1, -0.9)}}},
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
		grayt.Square{V1: LTF, V2: LTB, V3: RTB, V4: RTF},
		grayt.Square{V1: LBF, V2: RBF, V3: RTF, V4: LTF},
		grayt.Square{V1: LBB, V2: RBB, V3: RTB, V4: LTB},
		grayt.Square{V1: LBF, V2: LBB, V3: LTB, V4: LTF},
		grayt.Square{V1: RBF, V2: RBB, V3: RTB, V4: RTF},
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
		grayt.Square{V1: LTF, V2: LTB, V3: RTB, V4: RTF},
		grayt.Square{V1: LBF, V2: RBF, V3: RTF, V4: LTF},
		grayt.Square{V1: LBB, V2: RBB, V3: RTB, V4: LTB},
		grayt.Square{V1: LBF, V2: LBB, V3: LTB, V4: LTF},
		grayt.Square{V1: RBF, V2: RBB, V3: RTB, V4: RTF},
	}
}
