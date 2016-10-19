package main

import (
	"math"

	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("cornellbox", Scene{
		Camera: Cam(1.3),
		Triangles: JoinTriangles(
			JoinTriangles(
				ShortBlock(),
				TallBlock(),
				Floor,
				Ceiling,
				BackWall,
			).SetColour(White),
			LeftWall.SetColour(Red),
			RightWall.SetColour(Green),
			CeilingLight(),
		),
	})
}

func Cam(d float64) Camera {
	return Camera{
		Location:             Vect(0.5, 0.5, d),
		ViewDirection:        Vect(0.0, 0.0, -1.0),
		UpDirection:          Vect(0.0, 1.0, 0.0),
		FieldOfViewInDegrees: 2 * math.Asin(0.5/math.Sqrt(0.25+d*d)) * 180 / math.Pi,
		FocalLength:          0.5 + d,
		FocalRatio:           math.Inf(+1),
	}
}

var (
	Floor     = AlignedSquare(Vect(0, 0, 0), Vect(1, 0, -1))
	Ceiling   = AlignedSquare(Vect(0, 1, 0), Vect(1, 1, -1))
	BackWall  = AlignedSquare(Vect(0, 0, -1), Vect(1, 1, -1))
	LeftWall  = AlignedSquare(Vect(0, 0, 0), Vect(0, 1, -1))
	RightWall = AlignedSquare(Vect(1, 0, 0), Vect(1, 1, -1))
)

func CeilingLight() TriangleList {
	const size = 0.9
	return AlignedBox(
		Vect(size, 1.0, -size),
		Vect(1.0-size, 0.999, -1.0+size),
	).SetColour(White).SetEmittance(5.0)
}

func ShortBlock() TriangleList {
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
	)
}

func TallBlock() TriangleList {
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
	)
}
