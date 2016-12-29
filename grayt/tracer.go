package grayt

import (
	"image"
	"math/rand"
	"sync/atomic"
)

var rng = rand.New(rand.NewSource(1))

func traceImage(pxWide, pxHigh int, accel accelerationStructure, cam camera, quality int, completed *uint64) image.Image {

	accum := newAccumulator(pxWide, pxHigh)

	// Trace the image.
	pxPitch := 2.0 / float64(pxWide)
	for i := 0; i < quality; i++ {
		for pxY := 0; pxY < pxHigh; pxY++ {
			for pxX := 0; pxX < pxWide; pxX++ {
				x := (float64(pxX-pxWide/2) + rng.Float64()) * pxPitch
				y := (float64(pxY-pxHigh/2) + rng.Float64()) * pxPitch * -1.0
				r := cam.makeRay(x, y)
				r.dir = r.dir.Unit()
				accum.add(pxX, pxY, tracePath(accel, r))
				atomic.AddUint64(completed, 1)
			}
		}
	}

	return accum.toImage(1.0)
}

func tracePath(accel accelerationStructure, r ray) Colour {

	intersection, material, hit := accel.closestHit(r)
	if !hit {
		return Colour{0, 0, 0}
	}

	// Calculate probability of emitting.
	pEmit := 0.1
	if material.emittance != 0 {
		pEmit = 1.0
	}

	// Handle emit case.
	if rng.Float64() < pEmit {
		return material.colour.scale(material.emittance / pEmit)
	}

	// Find where the ray hit. Reduce the intersection distance by a small
	// amount so that reflected rays don't intersect with it immediately.
	hitLoc := r.at(addULPs(intersection.distance, -50))

	// Orient the unit normal towards the ray origin.
	if intersection.unitNormal.dot(r.dir) > 0 {
		intersection.unitNormal = intersection.unitNormal.Scale(-1.0)
	}

	// Create a random vector on the hemisphere towards the normal.
	rnd := Vector{rng.NormFloat64(), rng.NormFloat64(), rng.NormFloat64()}
	rnd = rnd.Unit()
	if rnd.dot(intersection.unitNormal) < 0 {
		rnd = rnd.Scale(-1.0)
	}

	// Apply the BRDF (bidirectional reflection distribution function).
	brdf := rnd.dot(intersection.unitNormal)

	return tracePath(accel, ray{start: hitLoc, dir: rnd}).
		scale(brdf / (1 - pEmit)).
		mul(material.colour)
}
