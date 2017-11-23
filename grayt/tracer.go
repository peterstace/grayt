package grayt

import (
	"image"
	"log"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

func TraceImage(pxWide int, scene Scene, quality, numWorkers int, accum *accumulator, completed *uint64) image.Image {
	pxHigh := scene.Camera.pxHigh(pxWide)
	cam := newCamera(scene.Camera)
	accel := newGrid(4, scene.Objects)

	finished := make(chan *pixelGrid)
	gridPool := make(chan *pixelGrid, numWorkers)
	for i := 0; i < numWorkers; i++ {
		gridPool <- &pixelGrid{
			pixels: make([]Colour, pxWide*pxHigh),
			wide:   pxWide,
			high:   pxHigh,
		}
	}

	go func() {
		pxPitch := 2.0 / float64(pxWide)
		for q := accum.count; q < quality; q++ {
			go func(q int, grid *pixelGrid) {
				tr := tracer{
					accel: accel,
					sky:   scene.Sky,
					rng:   rand.New(rand.NewSource(int64(q))),
				}
				for pxY := 0; pxY < pxHigh; pxY++ {
					for pxX := 0; pxX < pxWide; pxX++ {
						x := (float64(pxX-pxWide/2) + tr.rng.Float64()) * pxPitch
						y := (float64(pxY-pxHigh/2) + tr.rng.Float64()) * pxPitch * -1.0
						r := cam.makeRay(x, y, tr.rng)
						r.dir = r.dir.Unit()
						var c Colour
						if !*normals {
							c = tr.tracePath(r)
						} else {
							c = tr.traceNormal(r)
						}
						grid.set(pxX, pxY, c)
						atomic.AddUint64(completed, 1)
					}
				}
				finished <- grid
			}(q, <-gridPool)
		}
	}()

	ticker := time.NewTicker(time.Second)
	for q := accum.count; q < quality; {
		select {
		case grid := <-finished:
			accum.merge(grid)
			gridPool <- grid
			q++
		case <-ticker.C:
			if err := accum.save(); err != nil {
				log.Fatal("could not save snapshot:", err)
			}
		}
	}
	return accum.toImage(1.0)
}

type tracer struct {
	sky   func(Vector) Colour
	accel accelerationStructure
	rng   *rand.Rand
}

func (t *tracer) tracePath(r ray) Colour {
	assertUnit(r.dir)
	intersection, material, hit := t.accel.closestHit(r)
	if !hit {
		if t.sky == nil {
			return Colour{0, 0, 0}
		}
		assertUnit(r.dir)
		return t.sky(r.dir)
	}
	assertUnit(intersection.unitNormal)

	// Calculate probability of emitting.
	pEmit := 0.1
	if material.emittance != 0 {
		pEmit = 1.0
	}

	// Handle emit case.
	if t.rng.Float64() < pEmit {
		return material.colour.scale(material.emittance / pEmit)
	}

	offsetScale := -math.Copysign(addULPs(1.0, 1e5)-1.0, r.dir.Dot(intersection.unitNormal))
	offset := intersection.unitNormal.Scale(offsetScale)
	hitLoc := r.at(intersection.distance).Add(offset)

	// Orient the unit normal towards the ray origin.
	if intersection.unitNormal.Dot(r.dir) > 0 {
		intersection.unitNormal = intersection.unitNormal.Scale(-1.0)
	}

	if material.mirror {

		reflected := r.dir.Sub(intersection.unitNormal.Scale(2 * intersection.unitNormal.Dot(r.dir)))
		return t.tracePath(ray{start: hitLoc, dir: reflected})

	} else {

		// Create a random vector on the hemisphere towards the normal.
		rnd := Vector{t.rng.NormFloat64(), t.rng.NormFloat64(), t.rng.NormFloat64()}
		rnd = rnd.Unit()
		if rnd.Dot(intersection.unitNormal) < 0 {
			rnd = rnd.Scale(-1.0)
		}

		// Apply the BRDF (bidirectional reflection distribution function).
		brdf := rnd.Dot(intersection.unitNormal)

		return t.tracePath(ray{start: hitLoc, dir: rnd}).
			scale(brdf / (1 - pEmit)).
			mul(material.colour)
	}
}

func (t *tracer) traceNormal(r ray) Colour {
	intersection, _, hit := t.accel.closestHit(r)
	if !hit {
		return Colour{}
	}
	norm := intersection.unitNormal.Add(Vect(1, 1, 1)).Scale(0.5)
	return Colour{norm.X, norm.Y, norm.Z}
}
