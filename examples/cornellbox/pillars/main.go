package main

import (
	"math/rand"

	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("pillars", Scene{
		Camera: Cam(1.3),
		Objects: Group(
			Floor,
			Ceiling,
			BackWall,
			Group(
				pillar(Vect(0.30, 0.8, -0.75), 0.18),
				pillar(Vect(0.69, 0.7, -0.75), 0.17),
				pillar(Vect(0.17, 0.3, -0.25), 0.15),
				pillar(Vect(0.82, 0.4, -0.25), 0.16),
			).With(
				ColourRGB(0x80c0ff),
			),
			LeftWall.With(ColourRGB(Red)),
			RightWall.With(ColourRGB(Green)),
			Group(
				CeilingLight(),
				AlignedBox(
					Vect(0.96, 0.0, -0.04),
					Vect(0.70, 0.001, -0.06),
				),
				AlignedBox(
					Vect(0.04, 0.0, -0.04),
					Vect(0.30, 0.001, -0.06),
				),
			).With(
				Emittance(5.0),
			),
		),
	})
}

var rnd = rand.New(rand.NewSource(0))

func pillar(top Vector, radius float64) ObjectList {

	var objs ObjectList

	const plyHeight = 0.05
	segments := int(top.Y/plyHeight + 0.5)
	for i := 0; i < segments; i++ {
		r := radius * (1.0 - rnd.Float64()*0.5)
		objs = append(objs, AlignedBox(
			Vect(top.X-r, float64(i)*plyHeight, top.Z+r),
			Vect(top.X+r, (float64(i)+1)*plyHeight, top.Z-r),
		)...)
	}
	return objs
}
