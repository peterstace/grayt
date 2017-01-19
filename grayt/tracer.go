package grayt

import (
	"image"
	"math/rand"
	"sync"
	"sync/atomic"
)

func TraceImage(pxWide int, scene Scene, quality, numWorkers int, completed *uint64) image.Image {

	pxHigh := scene.Camera.pxHigh(pxWide)

	cam := newCamera(scene.Camera)
	accum := newAccumulator(pxWide, pxHigh)
	accel := newGrid(4, scene.Objects)

	var wg sync.WaitGroup
	wg.Add(quality)

	sem := make(chan struct{}, numWorkers)

	pxPitch := 2.0 / float64(pxWide)
	for q := 0; q < quality; q++ {
		sem <- struct{}{}
		go func(q int) {
			rng := rand.New(rand.NewSource(int64(q)))
			for pxY := 0; pxY < pxHigh; pxY++ {
				for pxX := 0; pxX < pxWide; pxX++ {
					x := (float64(pxX-pxWide/2) + rng.Float64()) * pxPitch
					y := (float64(pxY-pxHigh/2) + rng.Float64()) * pxPitch * -1.0
					r := cam.makeRay(x, y, rng)
					r.dir = r.dir.Unit()
					accum.add(pxX, pxY, tracePath(accel, r, rng), q)
					atomic.AddUint64(completed, 1)
				}
			}
			<-sem
			wg.Done()
		}(q)
	}
	wg.Wait()

	return createImage(accum)
}

func tracePath(accel accelerationStructure, r ray, rng *rand.Rand) Colour {

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
	hitLoc := r.at(addULPs(intersection.distance, -ulpFudgeFactor))

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

	return tracePath(accel, ray{start: hitLoc, dir: rnd}, rng).
		scale(brdf / (1 - pEmit)).
		mul(material.colour)
}
