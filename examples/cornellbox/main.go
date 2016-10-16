package main

import (
	"math"

	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("cornellbox", Scene{
		Camera:    cam(),
		Triangles: tris(),
	})
}

func cam() Camera {
	const D = 1.3 // Estimated.
	return Camera{
		Location:             Vect(0.5, 0.5, D),
		ViewDirection:        Vect(0.0, 0.0, -1.0),
		UpDirection:          Vect(0.0, 1.0, 0.0),
		FieldOfViewInDegrees: 2 * math.Asin(0.5/math.Sqrt(0.25+D*D)) * 180 / math.Pi,
		FocalLength:          0.5 + D,
		FocalRatio:           math.Inf(+1),
	}
}

func tris() []Triangle {
	const size = 0.9
	ts := alignedBox(
		Vect(size, 1.0, -size),
		Vect(1.0-size, 0.999, -1.0+size),
		5.0,
		White,
	)
	ts = append(ts, box()...)
	ts = append(ts, shortBlock()...)
	ts = append(ts, tallBlock()...)
	return ts
}

func alignedBox(a, b Vector, emittance float64, colour Colour) []Triangle {

	a1 := Vect(b.X, a.Y, a.Z)
	a2 := Vect(a.X, b.Y, a.Z)
	a3 := Vect(a.X, a.Y, b.Z)
	b1 := Vect(a.X, b.Y, b.Z)
	b2 := Vect(b.X, a.Y, b.Z)
	b3 := Vect(b.X, b.Y, a.Z)

	ts := []Triangle{

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

func square(a, b, c, d Vector, emittance float64, colour Colour) []Triangle {
	return []Triangle{
		{A: a, B: b, C: c, Emittance: emittance, Colour: colour},
		{A: c, B: d, C: a, Emittance: emittance, Colour: colour},
	}
}

func alignedSquare(a, b Vector, emittance float64, colour Colour) []Triangle {
	var c, d Vector
	switch {
	case a.X == b.X:
		c = Vector{X: a.X, Y: a.Y, Z: b.Z}
		d = Vector{X: a.X, Y: b.Y, Z: a.Z}
	case a.Y == b.Y:
		c = Vector{X: a.X, Y: a.Y, Z: b.Z}
		d = Vector{X: b.X, Y: a.Y, Z: a.Z}
	case a.Z == b.Z:
		c = Vector{X: a.X, Y: b.Y, Z: a.Z}
		d = Vector{X: b.X, Y: a.Y, Z: a.Z}
	default:
		panic("a and b line in a common aligned plane")

	}
	return []Triangle{
		{A: a, B: c, C: d, Colour: colour, Emittance: emittance},
		{A: b, B: c, C: d, Colour: colour, Emittance: emittance},
	}
}

func box() []Triangle {
	var ts []Triangle
	ts = append(ts, alignedSquare(Vect(0, 0, 0), Vect(1, 0, -1), 0, White)...)
	ts = append(ts, alignedSquare(Vect(0, 1, 0), Vect(1, 1, -1), 0, White)...)
	ts = append(ts, alignedSquare(Vect(0, 0, 0), Vect(0, 1, -1), 0, Red)...)
	ts = append(ts, alignedSquare(Vect(1, 0, 0), Vect(1, 1, -1), 0, Green)...)
	ts = append(ts, alignedSquare(Vect(0, 0, -1), Vect(1, 1, -1), 0, White)...)
	return ts
}

func shortBlock() []Triangle {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = Vect(0.76, 0.00, -0.12)
		LBB = Vect(0.85, 0.00, -0.41)
		RBF = Vect(0.47, 0.00, -0.21)
		RBB = Vect(0.56, 0.00, -0.49)
		LTF = Vect(0.76, 0.30, -0.12)
		LTB = Vect(0.85, 0.30, -0.41)
		RTF = Vect(0.47, 0.30, -0.21)
		RTB = Vect(0.56, 0.30, -0.49)
	)
	var ts []Triangle
	ts = append(ts, square(LTF, LTB, RTB, RTF, 0, White)...)
	ts = append(ts, square(LBF, RBF, RTF, LTF, 0, White)...)
	ts = append(ts, square(LBB, RBB, RTB, LTB, 0, White)...)
	ts = append(ts, square(LBF, LBB, LTB, LTF, 0, White)...)
	ts = append(ts, square(RBF, RBB, RTB, RTF, 0, White)...)
	return ts
}

func tallBlock() []Triangle {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = Vect(0.52, 0.00, -0.54)
		LBB = Vect(0.43, 0.00, -0.83)
		RBF = Vect(0.23, 0.00, -0.45)
		RBB = Vect(0.14, 0.00, -0.74)
		LTF = Vect(0.52, 0.60, -0.54)
		LTB = Vect(0.43, 0.60, -0.83)
		RTF = Vect(0.23, 0.60, -0.45)
		RTB = Vect(0.14, 0.60, -0.74)
	)
	var ts []Triangle
	ts = append(ts, square(LTF, LTB, RTB, RTF, 0, White)...)
	ts = append(ts, square(LBF, RBF, RTF, LTF, 0, White)...)
	ts = append(ts, square(LBB, RBB, RTB, LTB, 0, White)...)
	ts = append(ts, square(LBF, LBB, LTB, LTF, 0, White)...)
	ts = append(ts, square(RBF, RBB, RTB, RTF, 0, White)...)
	return ts
}
