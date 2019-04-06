package cornellbox

import (
	"math"

	"github.com/peterstace/grayt/scene"
	"github.com/peterstace/grayt/scene/dsl"
	"github.com/peterstace/grayt/xmath"
)

func CornellCam(d float64) scene.Camera {
	c := dsl.DefaultCamera()
	c.Location = xmath.Vect(0.5, 0.5, d)
	c.LookingAt = xmath.Vect(0.5, 0.5, -0.5)
	c.FieldOfViewInRadians = 2 * math.Asin(0.5/math.Sqrt(0.25+d*d))
	return c
}

var (
	CornellFloor     = dsl.AlignedSquare(xmath.Vect(0, 0, 0), xmath.Vect(1, 0, -1))
	CornellCeiling   = dsl.AlignedSquare(xmath.Vect(0, 1, 0), xmath.Vect(1, 1, -1))
	CornellBackWall  = dsl.AlignedSquare(xmath.Vect(0, 0, -1), xmath.Vect(1, 1, -1))
	CornellLeftWall  = dsl.AlignedSquare(xmath.Vect(0, 0, 0), xmath.Vect(0, 1, -1))
	CornellRightWall = dsl.AlignedSquare(xmath.Vect(1, 0, 0), xmath.Vect(1, 1, -1))
)

func CornellCeilingLight() scene.Surface {
	const size = 0.9
	return scene.Surface{
		AlignedBoxes: []scene.AlignedBox{{
			CornerA: xmath.Vect(size, 1.0, -size),
			CornerB: xmath.Vect(1.0-size, 0.999, -1.0+size),
		}},
	}
}

func CornellShortBlock() scene.Surface {
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
	return dsl.MergeSurfaces(
		dsl.Square(LTF, LTB, RTB, RTF),
		dsl.Square(LBF, RBF, RTF, LTF),
		dsl.Square(LBB, RBB, RTB, LTB),
		dsl.Square(LBF, LBB, LTB, LTF),
		dsl.Square(RBF, RBB, RTB, RTF),
	)
}

func CornellTallBlock() scene.Surface {
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
	return dsl.MergeSurfaces(
		dsl.Square(LTF, LTB, RTB, RTF),
		dsl.Square(LBF, RBF, RTF, LTF),
		dsl.Square(LBB, RBB, RTB, LTB),
		dsl.Square(LBF, LBB, LTB, LTF),
		dsl.Square(RBF, RBB, RTB, RTF),
	)
}
