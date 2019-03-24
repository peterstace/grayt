package grayt

import (
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/trace"
)

type render struct {
	completed int64
	traceRate int64

	requestedWorkers int64
	actualWorkers    int64

	// Static configuration.
	// TODO: Should ensure that these are not modified once the render is started.
	scene trace.Scene
	accum *accumulator
}

func newRender(scene trace.Scene, acc *accumulator) *render {
	return &render{
		completed: acc.passes * int64(acc.wide) * int64(acc.high),
		scene:     scene,
		accum:     acc,
	}
}

type status struct {
	completed        int64
	passes           int64
	traceRate        int64
	requestedWorkers int64
	actualWorkers    int64
}

func (r *render) status() status {
	return status{
		completed:        atomic.LoadInt64(&r.completed),
		passes:           atomic.LoadInt64(&r.accum.passes),
		traceRate:        atomic.LoadInt64(&r.traceRate),
		requestedWorkers: atomic.LoadInt64(&r.requestedWorkers),
		actualWorkers:    atomic.LoadInt64(&r.actualWorkers),
	}
}

func (r *render) setWorkers(workers int64) {
	atomic.StoreInt64(&r.requestedWorkers, workers)
}

func (r *render) traceImage() {
	accel := trace.NewGrid(4, r.scene.Objects)

	finished := make(chan *pixelGrid)
	gridPool := make(chan *pixelGrid)

	// Monitor trace rate.
	go func() {
		var lastCompleted int64
		const samplePeriod = 5 * time.Second
		ticker := time.NewTicker(samplePeriod)
		for {
			<-ticker.C
			completed := atomic.LoadInt64(&r.completed)
			if lastCompleted != 0 {
				sample := (completed - lastCompleted) * int64(time.Second) / int64(samplePeriod)
				atomic.StoreInt64(&r.traceRate, sample)
			}
			lastCompleted = completed
		}
	}()

	// Control size of worker pool.
	go func() {
		var dispatchedWorkers int64
		for {
			time.Sleep(100 * time.Millisecond)
			for dispatchedWorkers < atomic.LoadInt64(&r.requestedWorkers) {
				log.Println("dispatching worker")
				dispatchedWorkers++
				gridPool <- &pixelGrid{
					pixels: make([]colour.Colour, r.accum.wide*r.accum.high),
					wide:   r.accum.wide,
					high:   r.accum.high,
				}
			}
			for dispatchedWorkers > atomic.LoadInt64(&r.requestedWorkers) {
				log.Println("cancelling worker")
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
		pxPitch := 2.0 / float64(r.accum.wide)
		for i := 0; true; i++ {
			// TODO: Could pull off the worker pool at this point, rather than in separate goroutine.
			go func(i int, grid *pixelGrid) {
				atomic.AddInt64(&r.actualWorkers, 1)
				rng := rand.New(rand.NewSource(int64(i)))
				tr := trace.NewTracer(accel, rng)
				for pxY := 0; pxY < r.accum.high; pxY++ {
					for pxX := 0; pxX < r.accum.wide; pxX++ {
						x := (float64(pxX-r.accum.wide/2) + rng.Float64()) * pxPitch
						y := (float64(pxY-r.accum.high/2) + rng.Float64()) * pxPitch * -1.0
						cr := r.scene.Camera.MakeRay(x, y, rng)
						cr.Dir = cr.Dir.Unit()
						c := tr.TracePath(cr)
						grid.set(pxX, pxY, c)
						atomic.AddInt64(&r.completed, 1)
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
		gridPool <- grid
	}
}
