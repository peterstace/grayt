package main

import (
	"math"
	"math/rand"

	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("pillars", Scene{
		Camera: Cam(1.3),
		Triangles: JoinTriangles(
			JoinTriangles(
				Floor,
				Ceiling,
				BackWall,
			).SetColour(White),
			JoinTriangles(
				pillar(Vect(0.30, 0.8, -0.75), 0.18),
				pillar(Vect(0.69, 0.7, -0.75), 0.17),
				pillar(Vect(0.17, 0.3, -0.25), 0.15),
				pillar(Vect(0.82, 0.4, -0.25), 0.16),
			).SetColour(Colour{R: 0.5, G: 0.75, B: 1.0}),
			LeftWall.SetColour(Red),
			RightWall.SetColour(Green),
			CeilingLight(),
			AlignedBox(
				Vect(0.04, 0.0, -0.04),
				Vect(0.30, 0.001, -0.06),
			).SetColour(White).SetEmittance(5.0),
			AlignedBox(
				Vect(0.96, 0.0, -0.04),
				Vect(0.70, 0.001, -0.06),
			).SetColour(White).SetEmittance(5.0),
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

var rnd = rand.New(rand.NewSource(0))

func pillar(top Vector, radius float64) TriangleList {

	var tris TriangleList

	const plyHeight = 0.05
	segments := int(top.Y/plyHeight + 0.5)
	for i := 0; i < segments; i++ {
		r := radius * (1.0 - rnd.Float64()*0.5)
		tris = JoinTriangles(tris, AlignedBox(
			Vect(top.X-r, float64(i)*plyHeight, top.Z+r),
			Vect(top.X+r, (float64(i)+1)*plyHeight, top.Z-r),
		))
	}
	return tris
}
