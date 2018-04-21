package grayt

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

type render struct {
	completed uint64
	passes    uint64

	requestedWorkers int64
	actualWorkers    int64

	// Static configuration.
	// TODO: Should ensure that these are not modified once the render is started.
	pxWide int
	scene  Scene
	accum  *accumulator
}

func newRender(pxWide int, scene Scene, acc *accumulator) *render {
	return &render{
		pxWide: pxWide,
		scene:  scene,
		accum:  acc,
	}
}

type status struct {
	completed        uint64
	passes           uint64
	requestedWorkers int64
	actualWorkers    int64
}

func (r *render) status() status {
	return status{
		completed:        atomic.LoadUint64(&r.completed),
		passes:           atomic.LoadUint64(&r.passes),
		requestedWorkers: atomic.LoadInt64(&r.requestedWorkers),
		actualWorkers:    atomic.LoadInt64(&r.actualWorkers),
	}
}

func (r *render) setWorkers(workers int64) {
	atomic.StoreInt64(&r.requestedWorkers, workers)
}

func (r *render) traceImage() {
	pxHigh := r.scene.Camera.pxHigh(r.pxWide)
	cam := newCamera(r.scene.Camera)
	accel := newGrid(4, r.scene.Objects)

	finished := make(chan *pixelGrid)
	gridPool := make(chan *pixelGrid)

	go func() {
		var dispatchedWorkers int64
		for {
			time.Sleep(time.Second)
			for dispatchedWorkers < atomic.LoadInt64(&r.requestedWorkers) {
				dispatchedWorkers++
				gridPool <- &pixelGrid{
					pixels: make([]Colour, r.pxWide*pxHigh),
					wide:   r.pxWide,
					high:   pxHigh,
				}
			}
			for dispatchedWorkers > atomic.LoadInt64(&r.requestedWorkers) {
				dispatchedWorkers--
				// Run in goroutine, since we can't pull off the queue until a
				// pass finishes. We might not even pull off the next available
				// pixel grid, since the worker launcher might pick up the
				// spare grid first.
				go func() { <-gridPool }()
			}
		}
	}()

	// Launch workers.
	go func() {
		pxPitch := 2.0 / float64(r.pxWide)
		for i := 0; true; i++ {
			go func(i int, grid *pixelGrid) {
				atomic.AddInt64(&r.actualWorkers, 1)
				tr := tracer{
					accel: accel,
					sky:   r.scene.Sky,
					rng:   rand.New(rand.NewSource(int64(i))),
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
				atomic.AddInt64(&r.actualWorkers, -1)
				finished <- grid
			}(i, <-gridPool)
		}
	}()

	// Coordination point for merging worker results.
	for grid := range finished {
		r.accum.merge(grid)
		atomic.AddUint64(&r.passes, 1)
		gridPool <- grid
	}
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
