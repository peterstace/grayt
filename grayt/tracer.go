package grayt

import (
	"image"
	"log"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

type render struct {
	completed uint64

	// Static configuration.
	// TODO: Should ensure that these are not modified once the render is started.
	pxWide     int
	quality    int
	numWorkers int
	scene      Scene
}

func (r *render) traceImage(accum *accumulator) image.Image {
	pxHigh := r.scene.Camera.pxHigh(r.pxWide)
	cam := newCamera(r.scene.Camera)
	accel := newGrid(4, r.scene.Objects)

	finished := make(chan *pixelGrid)
	gridPool := make(chan *pixelGrid, r.numWorkers)
	for i := 0; i < r.numWorkers; i++ {
		gridPool <- &pixelGrid{
			pixels: make([]Colour, r.pxWide*pxHigh),
			wide:   r.pxWide,
			high:   pxHigh,
		}
	}

	go func() {
		pxPitch := 2.0 / float64(r.pxWide)
		for q := accum.count; q < r.quality; q++ {
			go func(q int, grid *pixelGrid) {
				tr := tracer{
					accel: accel,
					sky:   r.scene.Sky,
					rng:   rand.New(rand.NewSource(int64(q))),
				}
				for pxY := 0; pxY < pxHigh; pxY++ {
					for pxX := 0; pxX < r.pxWide; pxX++ {
						x := (float64(pxX-r.pxWide/2) + tr.rng.Float64()) * pxPitch
						y := (float64(pxY-pxHigh/2) + tr.rng.Float64()) * pxPitch * -1.0
						cr := cam.makeRay(x, y, tr.rng)
						cr.dir = cr.dir.Unit()
						c := tr.tracePath(cr)
						grid.set(pxX, pxY, c)
						atomic.AddUint64(&r.completed, 1)
					}
				}
				finished <- grid
			}(q, <-gridPool)
		}
	}()

	ticker := time.NewTicker(10 * time.Minute)
	for q := accum.count; q < r.quality; {
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
