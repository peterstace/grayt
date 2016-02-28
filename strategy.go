package grayt

import (
	"image"
	"log"
	"math/rand"
	"time"
)

type strategy struct {
}

func (s *strategy) traceImage(pxHigh, pxWide int, scene Scene, quality int) image.Image {

	acc := newAccumulator(pxHigh, pxWide)

	done := make(chan struct{})
	go func() {
		final := false
		for {
			select {
			case <-done:
				final = true
			default:
			}
			log.Printf("%4.1f%%\n", float64(acc.getTotal())/float64(quality*pxWide*pxHigh)*100)
			if final {
				log.Print("Done.")
				done <- struct{}{}
				return
			}
			time.Sleep(time.Second)
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
