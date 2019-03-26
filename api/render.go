package api

import (
	"image"
	"image/color"
	"log"
	"sync"
	"time"
)

type render struct {
	scene   string
	pxWide  int
	pxHigh  int
	created time.Time

	cnd            *sync.Cond
	desiredWorkers int
}

func (r *render) work() {
	for {
		log.Println("...working...")
		time.Sleep(time.Second)
	}
}

func (r *render) image() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, r.pxWide, r.pxHigh))
	for x := 0; x < r.pxWide; x++ {
		for y := 0; y < r.pxHigh; y++ {
			img.Set(x, y, color.Gray{0xff})
		}
	}
	return img
}
