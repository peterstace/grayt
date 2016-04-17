package grayt

import (
	"image"
	"math/rand"
	"sync/atomic"
	"time"
)

type strategy struct {
}

func (s *strategy) traceImage(pxHigh, pxWide int, scene Scene, quality int) image.Image {

	acc := newAccumulator(pxHigh, pxWide)

	// MUST only be used atomically, since it's shared by 2 goroutines - the
	// tracing goroutine and the CLI goroutine.
	var completed uint64

	// Update the CLI (in a separate goroutine).
	cli := newCLI()
	done := make(chan struct{})
	go func() {
		total := uint64(pxWide * pxHigh * quality)
		for {
			var exit bool
			select {
			case <-done:
				exit = true
			case <-time.After(cliUpdatePeriod):
			}
			cli.update(atomic.LoadUint64(&completed), total)
			if exit {
				cli.done()
				done <- struct{}{}
				return
			}
		}
	}()

	// Trace the image.
	w := newWorld(scene.Entities)
	pxPitch := 2.0 / float64(pxWide)
	for pxX := 0; pxX < pxWide; pxX++ {
		for pxY := 0; pxY < pxHigh; pxY++ {
			for i := 0; i < quality; i++ {
				x := (float64(pxX-pxWide/2) + rand.Float64()) * pxPitch
				y := (float64(pxY-pxHigh/2) + rand.Float64()) * pxPitch * -1.0
				r := scene.Camera.MakeRay(x, y)
				r.Dir = r.Dir.Unit()
				acc.add(pxX, pxY, tracePath(w, r))
				atomic.AddUint64(&completed, 1)
			}
		}
	}

	// Tell the CLI to finish, then wait for its reply.
	done <- struct{}{}
	<-done

	return acc.toImage(1.0)
}

// TODO: Also output some kind of meta file about the image that was generated?
