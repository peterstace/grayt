package main

import "github.com/peterstace/grayt/scene"

func newAccelerationStructure(ts []scene.Triangle) accelerationStructure {
	return accelerationStructure{ts}
}

type accelerationStructure struct {
	tris []scene.Triangle
}

func (a accelerationStructure) closestHit(r Ray) (intersection, bool) {
	var closest struct {
		intersection intersection
		hit          bool
	}
	for i := range a.tris {
		intersection, hit := w.entities[i].Surface.Intersect(r)
		if !hit {
			continue
		}
		if !closest.hit || intersection.Distance < closest.intersection.distance {
			closest.intersection = intersection
			closest.hit = true
		}
	}
	return closest.intersection, closest.hit
}
