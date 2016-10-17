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
	return JoinTriangles(
		AlignedBox(
			Vect(size, 1.0, -size),
			Vect(1.0-size, 0.999, -1.0+size),
		).SetColour(White).SetEmittance(5.0),
		box(),
		shortBlock(),
		tallBlock(),
	)
}

func box() TriangleList {
	return JoinTriangles(
		JoinTriangles(
			AlignedSquare(Vect(0, 0, 0), Vect(1, 0, -1)),
			AlignedSquare(Vect(0, 1, 0), Vect(1, 1, -1)),
			AlignedSquare(Vect(0, 0, -1), Vect(1, 1, -1)),
		).SetColour(White),
		AlignedSquare(Vect(0, 0, 0), Vect(0, 1, -1)).SetColour(Red),
		AlignedSquare(Vect(1, 0, 0), Vect(1, 1, -1)).SetColour(Green),
	)
}

func shortBlock() TriangleList {
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
	return JoinTriangles(
		Square(LTF, LTB, RTB, RTF),
		Square(LBF, RBF, RTF, LTF),
		Square(LBB, RBB, RTB, LTB),
		Square(LBF, LBB, LTB, LTF),
		Square(RBF, RBB, RTB, RTF),
	).SetColour(White)
}

func tallBlock() TriangleList {
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
	return JoinTriangles(
		Square(LTF, LTB, RTB, RTF),
		Square(LBF, RBF, RTF, LTF),
		Square(LBB, RBB, RTB, LTB),
		Square(LBF, LBB, LTB, LTF),
		Square(RBF, RBB, RTB, RTF),
	).SetColour(White)
}
