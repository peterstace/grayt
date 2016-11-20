package main

import (
	"math"
	"math/rand"

	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("sphere_tree", Scene{
		Camera: Cam(1.3),
		Objects: Group(
			Tree(),
			Floor,
			Ceiling,
			BackWall,
			LeftWall.With(ColourRGB(Red)),
			RightWall.With(ColourRGB(Green)),
			CeilingLight().With(Emittance(5.0)),
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
	BackWall  = ZPlane(-1)
	LeftWall  = AlignedSquare(Vect(0, 0, 0), Vect(0, 1, -1))
	RightWall = AlignedSquare(Vect(1, 0, 0), Vect(1, 1, -1))
)

func CeilingLight() ObjectList {
	const size = 0.9
	return AlignedBox(
		Vect(size, 1.0, -size),
		Vect(1.0-size, 0.999, -1.0+size),
	)
}

type sphere struct {
	c Vector
	r float64
}

func Tree() ObjectList {
	var ol ObjectList
	for _, s := range recurse(sphere{Vect(0.5, 0, -0.5), 0.2}, 8) {
		ol = Group(ol, Sphere(s.c, s.r))
	}
	return ol
}

const radiusScaleDown = 0.8

var rnd = rand.New(rand.NewSource(2))

func recurse(parent sphere, level int) []sphere {

	if level == 0 {
		return []sphere{parent}
	}

	var spheres []sphere
	for i := 0; i < 2; i++ {
		var child sphere
		matches := false
		for !matches {
			rndUnit := Vector{rnd.NormFloat64(), rnd.NormFloat64(), rnd.NormFloat64()}.Unit()
			child.c = parent.c.Add(rndUnit.Scale(parent.r))
			child.r = radiusScaleDown * parent.r
			matches = true &&
				child.c.X > child.r &&
				child.c.X < 1.0-child.r &&
				child.c.Y > child.r &&
				child.c.Y < 1.0-child.r &&
				child.c.Z < -child.r &&
				child.c.Z > -1.0+child.r
		}
		spheres = append(spheres, recurse(child, level-1)...)
	}
	return append(spheres, parent)
}
