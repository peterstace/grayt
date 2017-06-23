package main

import (
	"math"

	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("platonic", scene())
}

func cam() CameraBlueprint {
	const d = 1.3
	return Camera().With(
		Location(Vect(0.5, 0.5, d)),
		LookingAt(Vect(0.5, 0.5, -0.5)),
		FieldOfViewInRadians(2*math.Asin(0.5/math.Sqrt(0.25+d*d))),
	)
}

func tetrahedron() ObjectList {

	invSqrt2 := 1.0 / math.Sqrt2
	v1 := Vect(-1, 0, -invSqrt2)
	v2 := Vect(+1, 0, -invSqrt2)
	v3 := Vect(0, -1, +invSqrt2)
	v4 := Vect(0, +1, +invSqrt2)

	const scale = 0.1
	offset := Vect(0.2, 0.5, -0.5)

	v1 = v1.Scale(scale).Add(offset)
	v2 = v2.Scale(scale).Add(offset)
	v3 = v3.Scale(scale).Add(offset)
	v4 = v4.Scale(scale).Add(offset)

	return Group(
		Triangle(v1, v2, v3),
		Triangle(v1, v2, v4),
		Triangle(v1, v3, v4),
		Triangle(v2, v3, v4),
	)
}

func scene() Scene {

	var (
		Floor     = AlignedSquare(Vect(0, 0, 0), Vect(1, 0, -1))
		Ceiling   = AlignedSquare(Vect(0, 1, 0), Vect(1, 1, -1))
		BackWall  = AlignedSquare(Vect(0, 0, -1), Vect(1, 1, -1))
		LeftWall  = AlignedSquare(Vect(0, 0, 0), Vect(0, 1, -1))
		RightWall = AlignedSquare(Vect(1, 0, 0), Vect(1, 1, -1))

		size         = 0.9
		CeilingLight = AlignedBox(
			Vect(size, 1.0, -size),
			Vect(1.0-size, 0.999, -1.0+size),
		)
	)

	return Scene{
		Camera: cam(),
		Objects: Group(
			Floor,
			Ceiling,
			BackWall,
			LeftWall.With(ColourRGB(Red)),
			RightWall.With(ColourRGB(Green)),
			CeilingLight.With(Emittance(5.0)),
			tetrahedron(),
		),
	}
}
