package grayt

import (
	"fmt"
	"image"
	"math/rand"
	"time"
)

type strategy struct {
}

func (s *strategy) traceImage(pxHigh, pxWide int, scene Scene, quality int) image.Image {

	start := time.Now()

	acc := newAccumulator(pxHigh, pxWide)

	done := make(chan struct{})
	go func() {

		final := false

		var throughputSmoothed float64
		samples := 0
		now := time.Now()

		for {
			select {
			case <-done:
				final = true
			default:
			}

			var nowDelta time.Duration
			newNow := time.Now()
			nowDelta, now = newNow.Sub(now), newNow
			newSamples := acc.getTotal()
			var samplesDelta int
			samplesDelta, samples = newSamples-samples, newSamples
			throughput := float64(samplesDelta) / nowDelta.Seconds()
			const alpha = 0.1
			throughputSmoothed = throughputSmoothed*(1.0-alpha) + throughput*alpha

			stats{
				elapsed:    time.Nanosecond * time.Duration(time.Now().Sub(start).Nanoseconds()/1e7*1e7),
				throughput: throughputSmoothed,
			}.display()

			if final {
				fmt.Printf("\nDone.\n")
				done <- struct{}{}
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	w := newWorld(scene.Entities)
	pxPitch := 2.0 / float64(pxWide)
	for i := 0; i < quality; i++ {
		for pxX := 0; pxX < pxWide; pxX++ {
			for pxY := 0; pxY < pxHigh; pxY++ {
				x := (float64(pxX-pxWide/2) + rand.Float64()) * pxPitch
				y := (float64(pxY-pxHigh/2) + rand.Float64()) * pxPitch * -1.0
				r := scene.Camera.MakeRay(x, y)
				r.Dir = r.Dir.Unit()
				acc.add(pxX, pxY, tracePath(w, r))
			}
		}
	}

	done <- struct{}{}
	<-done

	return acc.toImage(1.0)
}

type stats struct {
	elapsed    time.Duration
	throughput float64
}

func (s stats) display() {
	fmt.Print("\x1b[1G") // Move to column 1.
	fmt.Print("\x1b[2K") // Clear line.
	fmt.Printf("%v %v", s.elapsed, s.throughput)
}
