package cornellbox

import (
	"math"

	. "github.com/peterstace/grayt/grayt"
	"github.com/peterstace/grayt/xmath"
)

func Cam(d float64) CameraBlueprint {
	return Camera().With(
		Location(xmath.Vect(0.5, 0.5, d)),
		LookingAt(xmath.Vect(0.5, 0.5, -0.5)),
		FieldOfViewInRadians(2*math.Asin(0.5/math.Sqrt(0.25+d*d))),
	)
}

var (
	Floor     = AlignedSquare(xmath.Vect(0, 0, 0), xmath.Vect(1, 0, -1))
	Ceiling   = AlignedSquare(xmath.Vect(0, 1, 0), xmath.Vect(1, 1, -1))
	BackWall  = AlignedSquare(xmath.Vect(0, 0, -1), xmath.Vect(1, 1, -1))
	LeftWall  = AlignedSquare(xmath.Vect(0, 0, 0), xmath.Vect(0, 1, -1))
	RightWall = AlignedSquare(xmath.Vect(1, 0, 0), xmath.Vect(1, 1, -1))
)

func CeilingLight() ObjectList {
	const size = 0.9
	return AlignedBox(
		xmath.Vect(size, 1.0, -size),
		xmath.Vect(1.0-size, 0.999, -1.0+size),
	)
}

func ShortBlock() ObjectList {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = xmath.Vect(0.76, 0.00, -0.12)
		LBB = xmath.Vect(0.85, 0.00, -0.41)
		RBF = xmath.Vect(0.47, 0.00, -0.21)
		RBB = xmath.Vect(0.56, 0.00, -0.49)
		LTF = xmath.Vect(0.76, 0.30, -0.12)
		LTB = xmath.Vect(0.85, 0.30, -0.41)
		RTF = xmath.Vect(0.47, 0.30, -0.21)
		RTB = xmath.Vect(0.56, 0.30, -0.49)
	)
	return Group(
		Square(LTF, LTB, RTB, RTF),
		Square(LBF, RBF, RTF, LTF),
		Square(LBB, RBB, RTB, LTB),
		Square(LBF, LBB, LTB, LTF),
		Square(RBF, RBB, RTB, RTF),
	)
}

func TallBlock() ObjectList {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = xmath.Vect(0.52, 0.00, -0.54)
		LBB = xmath.Vect(0.43, 0.00, -0.83)
		RBF = xmath.Vect(0.23, 0.00, -0.45)
		RBB = xmath.Vect(0.14, 0.00, -0.74)
		LTF = xmath.Vect(0.52, 0.60, -0.54)
		LTB = xmath.Vect(0.43, 0.60, -0.83)
		RTF = xmath.Vect(0.23, 0.60, -0.45)
		RTB = xmath.Vect(0.14, 0.60, -0.74)
	)
	return Group(
		Square(LTF, LTB, RTB, RTF),
		Square(LBF, RBF, RTF, LTF),
		Square(LBB, RBB, RTB, LTB),
		Square(LBF, LBB, LTB, LTF),
		Square(RBF, RBB, RTB, RTF),
	)
}
