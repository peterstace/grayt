package main

import (
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/peterstace/grayt/scene"
)

func main() {

	s := scene.Scene{
		Camera:    cam(),
		Triangles: tris(),
	}

	f, err := os.Create(filepath.Base(os.Args[0]) + ".bin")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.WriteTo(f); err != nil {
		log.Fatal(err)
	}
}

var (
	up    = scene.Vect(0.0, 1.0, 0.0)
	down  = scene.Vect(0.0, -1.0, 0.0)
	left  = scene.Vect(-1.0, 0.0, 0.0)
	right = scene.Vect(1.0, 0.0, 0.0)
	back  = scene.Vect(0.0, 0.0, 1.0)
	zero  = scene.Vect(0.0, 0.0, 0.0)
	one   = scene.Vect(1.0, 1.0, -1.0)
)

var (
	white = scene.Colour{1, 1, 1}
	green = scene.Colour{0, 1, 0}
	red   = scene.Colour{1, 0, 0}
)

func cam() scene.Camera {
	const D = 1.3 // Estimated.
	return scene.Camera{
		Location:             scene.Vect(0.5, 0.5, D),
		ViewDirection:        scene.Vect(0.0, 0.0, -1.0),
		UpDirection:          up,
		FieldOfViewInDegrees: 2 * math.Asin(0.5/math.Sqrt(0.25+D*D)) * 180 / math.Pi,
		FocalLength:          0.5 + D,
		FocalRatio:           math.Inf(+1),
	}
}

func tris() []scene.Triangle {
	const size = 0.4
	ts := alignedBox(
		scene.Vect(size, 1.0, -size),
		scene.Vect(1.0-size, 0.999, -1.0+size),
		5.0,
		white,
	)
	ts = append(ts, box()...)
	ts = append(ts, shortBlock()...)
	ts = append(ts, tallBlock()...)
	return ts
}

func alignedBox(a, b scene.Vector, emittance float64, colour scene.Colour) []scene.Triangle {

	a1 := scene.Vector{X: b.X, Y: a.Y, Z: a.Z}
	a2 := scene.Vector{X: a.X, Y: b.Y, Z: a.Z}
	a3 := scene.Vector{X: a.X, Y: a.Y, Z: b.Z}
	b1 := scene.Vector{X: a.X, Y: b.Y, Z: b.Z}
	b2 := scene.Vector{X: b.X, Y: a.Y, Z: b.Z}
	b3 := scene.Vector{X: b.X, Y: b.Y, Z: a.Z}

	ts := []scene.Triangle{

		{A: a, B: a1, C: a2},
		{A: a, B: a2, C: a3},
		{A: a, B: a3, C: a1},

		{A: b, B: b1, C: b2},
		{A: b, B: b2, C: b3},
		{A: b, B: b3, C: b1},

		{A: a1, B: b2, C: b3},
		{A: a2, B: b3, C: b1},
		{A: a3, B: b1, C: b2},

		{A: b1, B: a2, C: a3},
		{A: b2, B: a3, C: a1},
		{A: b3, B: a1, C: a2},
	}

	for i := range ts {
		ts[i].Colour = colour
		ts[i].Emittance = emittance
	}

	return ts
}

func square(a, b, c, d scene.Vector, emittance float64, colour scene.Colour) []scene.Triangle {
	return []scene.Triangle{
		{A: a, B: b, C: c, Emittance: emittance, Colour: colour},
		{A: c, B: d, C: a, Emittance: emittance, Colour: colour},
	}
}

func alignedSquare(a, b scene.Vector, emittance float64, colour scene.Colour) []scene.Triangle {
	var c, d scene.Vector
	switch {
	case a.X == b.X:
		c = scene.Vector{X: a.X, Y: a.Y, Z: b.Z}
		d = scene.Vector{X: a.X, Y: b.Y, Z: a.Z}
	case a.Y == b.Y:
		c = scene.Vector{X: a.X, Y: a.Y, Z: b.Z}
		d = scene.Vector{X: b.X, Y: a.Y, Z: a.Z}
	case a.Z == b.Z:
		c = scene.Vector{X: a.X, Y: b.Y, Z: a.Z}
		d = scene.Vector{X: b.X, Y: a.Y, Z: a.Z}
	default:
		panic("a and b line in a common aligned plane")

	}
	return []scene.Triangle{
		{A: a, B: c, C: d, Colour: colour, Emittance: emittance},
		{A: b, B: c, C: d, Colour: colour, Emittance: emittance},
	}
}

func box() []scene.Triangle {
	var ts []scene.Triangle
	ts = append(ts,
		alignedSquare(scene.Vect(0, 0, 0), scene.Vect(1, 0, -1), 0, white)...)
	ts = append(ts,
		alignedSquare(scene.Vect(0, 1, 0), scene.Vect(1, 1, -1), 0, white)...)
	ts = append(ts,
		alignedSquare(scene.Vect(0, 0, 0), scene.Vect(0, 1, -1), 0, red)...)
	ts = append(ts,
		alignedSquare(scene.Vect(1, 0, 0), scene.Vect(1, 1, -1), 0, green)...)
	ts = append(ts,
		alignedSquare(scene.Vect(0, 0, -1), scene.Vect(1, 1, -1), 0, white)...)
	return ts
}

func shortBlock() []scene.Triangle {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = scene.Vect(0.76, 0.00, -0.12)
		LBB = scene.Vect(0.85, 0.00, -0.41)
		RBF = scene.Vect(0.47, 0.00, -0.21)
		RBB = scene.Vect(0.56, 0.00, -0.49)
		LTF = scene.Vect(0.76, 0.30, -0.12)
		LTB = scene.Vect(0.85, 0.30, -0.41)
		RTF = scene.Vect(0.47, 0.30, -0.21)
		RTB = scene.Vect(0.56, 0.30, -0.49)
	)
	var ts []scene.Triangle
	ts = append(ts, square(LTF, LTB, RTB, RTF, 0, white)...)
	ts = append(ts, square(LBF, RBF, RTF, LTF, 0, white)...)
	ts = append(ts, square(LBB, RBB, RTB, LTB, 0, white)...)
	ts = append(ts, square(LBF, LBB, LTB, LTF, 0, white)...)
	ts = append(ts, square(RBF, RBB, RTB, RTF, 0, white)...)
	return ts
}

func tallBlock() []scene.Triangle {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = scene.Vect(0.52, 0.00, -0.54)
		LBB = scene.Vect(0.43, 0.00, -0.83)
		RBF = scene.Vect(0.23, 0.00, -0.45)
		RBB = scene.Vect(0.14, 0.00, -0.74)
		LTF = scene.Vect(0.52, 0.60, -0.54)
		LTB = scene.Vect(0.43, 0.60, -0.83)
		RTF = scene.Vect(0.23, 0.60, -0.45)
		RTB = scene.Vect(0.14, 0.60, -0.74)
	)
	var ts []scene.Triangle
	ts = append(ts, square(LTF, LTB, RTB, RTF, 0, white)...)
	ts = append(ts, square(LBF, RBF, RTF, LTF, 0, white)...)
	ts = append(ts, square(LBB, RBB, RTB, LTB, 0, white)...)
	ts = append(ts, square(LBF, LBB, LTB, LTF, 0, white)...)
	ts = append(ts, square(RBF, RBB, RTB, RTF, 0, white)...)
	return ts
}
