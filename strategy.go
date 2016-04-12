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
	fmt.Printf(
		"Duration: %s Throughput: %s samples/sec",
		displayDuration(s.elapsed), displayFloat64(s.throughput),
	)
}

func displayFloat64(f float64) string {

	var thousands int

	for f >= 1000 {
		f /= 1000
		thousands++
	}

	suffix := [...]byte{' ', 'K', 'M', 'T', 'P', 'Y'}[thousands]

	if f < 10 {
		// 9.999K
		return fmt.Sprintf("%.3f%c", f, suffix)
	} else if f < 100 {
		// 99.99K
		return fmt.Sprintf("%.2f%c", f, suffix)
	} else if f < 1000 {
		// 999.9K
		return fmt.Sprintf("%.1f%c", f, suffix)
	}
	return fmt.Sprintf("%f", f)
}

func displayDuration(d time.Duration) string {
	h := d / time.Hour
	m := (d - h*time.Hour) / time.Minute
	s := (d - h*time.Hour - m*time.Minute) / time.Second
	return fmt.Sprintf(
		"%d%d:%d%d:%d%d",
		h/10, h%10, m/10, m%10, s/10, s%10,
	)
}
