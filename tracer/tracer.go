package main

import (
	"image"
	"math/rand"
)

func traceImage(pxWide, pxHigh int, accel accelerationStructure, cam camera) image.Image {

	accum := newAccumulator(pxWide, pxHigh)

	const quality = 10

	// Trace the image.
	pxPitch := 2.0 / float64(pxWide)
	for pxX := 0; pxX < pxWide; pxX++ {
		for pxY := 0; pxY < pxHigh; pxY++ {
			for i := 0; i < quality; i++ {
				x := (float64(pxX-pxWide/2) + rand.Float64()) * pxPitch
				y := (float64(pxY-pxHigh/2) + rand.Float64()) * pxPitch * -1.0
				r := cam.makeRay(x, y)
				r.dir = r.dir.unit()
				accum.add(pxX, pxY, tracePath(accel, r))
			}
		}
	}

	return accum.toImage(1.0)
}

func tracePath(accel accelerationStructure, r ray) colour {

	intersection, hit := accel.closestHit(r)
	if !hit {
		return colour{0, 0, 0}
	}

	// Calculate probability of emitting.
	pEmit := 0.1
	if material.Emittance != 0 {
		pEmit = 1.0
	}

	// Handle emit case.
	if rand.Float64() < pEmit {
		return material.colour.
			Scale(material.Emittance / pEmit)
	}

	// Find where the ray hit. Reduce the intersection distance by a small
	// amount so that reflected rays don't intersect with it immediately.
	hitLoc := r.At(addULPs(intersection.Distance, -50))

	// Orient the unit normal towards the ray origin.
	if intersection.UnitNormal.Dot(r.Dir) > 0 {
		intersection.UnitNormal = intersection.UnitNormal.Scale(-1.0)
	}

	// Create a random vector on the hemisphere towards the normal.
	rnd := Vector{rand.NormFloat64(), rand.NormFloat64(), rand.NormFloat64()}
	rnd = rnd.Unit()
	if rnd.Dot(intersection.UnitNormal) < 0 {
		rnd = rnd.Scale(-1.0)
	}

	// Apply the BRDF (bidirectional reflection distribution function).
	brdf := rnd.Dot(intersection.UnitNormal)

	return tracePath(w, Ray{Start: hitLoc, Dir: rnd}).
		Scale(brdf / (1 - pEmit)).
		Mul(material.Colour)
}
