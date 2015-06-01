package grayt

import (
	"math"
	"math/rand"
)

func trace(f func(Ray) Colour, c Camera, a Accumulator, spp int) {

	pxPitch := 2.0 / float64(a.wide)
	for pxX := 0; pxX < a.wide; pxX++ {
		for pxY := 0; pxY < a.high; pxY++ {

			//if pxX != 140 || pxY != 40 {
			//	continue
			//}

			for i := 0; i < spp; i++ {
				x := (float64(pxX-a.wide/2) + rand.Float64()) * pxPitch
				y := (float64(pxY-a.high/2) + rand.Float64()) * pxPitch * -1.0
				r := c.MakeRay(x, y)
				r.Dir = r.Dir.Unit()
				a.add(pxX, pxY, f(r))
			}
		}
	}
}

func closestHit(surfaces []Surface, r Ray) (Intersection, Surface) {
	var closest struct {
		intersection Intersection
		surface      Surface
	}
	for _, surface := range surfaces {
		intersection, hit := surface.Intersect(r)
		if !hit {
			continue
		}
		if closest.surface == nil || intersection.Distance < closest.intersection.Distance {
			closest.intersection = intersection
			closest.surface = surface
		}
	}
	return closest.intersection, closest.surface
}

func ulpDiff(a, b float64) uint64 {

	ulpA := math.Float64bits(a)
	ulpB := math.Float64bits(b)

	if ulpA > ulpB {
		return ulpA - ulpB
	} else {
		return ulpB - ulpA
	}
}
