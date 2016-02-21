package main

import (
	"math"

	"github.com/peterstace/grayt"
	"github.com/peterstace/grayt/geometry"
)

func main() {
	r := grayt.NewRunner()
	r.PxWide = 320
	r.PxHigh = 320
	r.Run(scene())
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
	white = grayt.Material{Colour: grayt.White, Emittance: 0.0}
	green = grayt.Material{Colour: grayt.Green, Emittance: 0.0}
	red   = grayt.Material{Colour: grayt.Red, Emittance: 0.0}
	blue  = grayt.Material{Colour: grayt.Blue, Emittance: 0.0}
)

func scene() grayt.Scene {
	ee := []grayt.Entity{{
		Material: grayt.Material{Colour: grayt.White, Emittance: 5},
		Surface:  geometry.Sphere(grayt.Vect(0.5, 1.0, -0.5), 0.25),
	}}
	ee = append(ee, box()...)
	for _, s := range tallBlock() {
		ee = append(ee, grayt.Entity{Material: white, Surface: s})
	}
	for _, s := range shortBlock() {
		ee = append(ee, grayt.Entity{Material: white, Surface: s})
	}
	return grayt.Scene{Camera: cam(), Entities: ee}
}

func cam() grayt.Camera {
	const D = 1.3 // Estimated.
	return grayt.NewRectilinearCamera(grayt.CameraConfig{
		Location:      grayt.Vect(0.5, 0.5, D),
		ViewDirection: grayt.Vect(0.0, 0.0, -1.0),
		UpDirection:   up,
		FieldOfView:   2 * math.Asin(0.5/math.Sqrt(0.25+D*D)),
		FocalLength:   0.5 + D,
		FocalRatio:    math.Inf(+1),
	})
}

func box() []grayt.Entity {
	return []grayt.Entity{
		{
			Material: white,
			Surface:  geometry.Plane(up, zero),
		},
		{
			Material: white,
			Surface:  geometry.Plane(down, one),
		},
		{
			Material: red,
			Surface:  geometry.Plane(right, zero),
		},
		{
			Material: green,
			Surface:  geometry.Plane(left, one),
		},
		{
			Material: white,
			Surface:  geometry.Plane(back, one),
		},
		{
			Material: white,
			Surface: geometry.AlignedBox(
				grayt.Vect(0.1, 0.0, -0.1),
				grayt.Vect(0.9, 0.1, -0.9),
			),
		},
	}
}

func shortBlock() []grayt.Surface {
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
	var ss []grayt.Surface
	ss = append(ss, geometry.Square(LTF, LTB, RTB, RTF)...)
	ss = append(ss, geometry.Square(LBF, RBF, RTF, LTF)...)
	ss = append(ss, geometry.Square(LBB, RBB, RTB, LTB)...)
	ss = append(ss, geometry.Square(LBF, LBB, LTB, LTF)...)
	ss = append(ss, geometry.Square(RBF, RBB, RTB, RTF)...)
	return ss
}

func tallBlock() []grayt.Surface {
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
	var ss []grayt.Surface
	ss = append(ss, geometry.Square(LTF, LTB, RTB, RTF)...)
	ss = append(ss, geometry.Square(LBF, RBF, RTF, LTF)...)
	ss = append(ss, geometry.Square(LBB, RBB, RTB, LTB)...)
	ss = append(ss, geometry.Square(LBF, LBB, LTB, LTF)...)
	ss = append(ss, geometry.Square(RBF, RBB, RTB, RTF)...)
	return ss
}
