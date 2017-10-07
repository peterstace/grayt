package engine

import (
	"fmt"
	"image"
)

type Status struct {
	// Accessed atomically
	Done  int64
	Total int64
}

func TraceImage(pxWide int, scene func(*API), status *Status) image.Image {
	api := newAPI()
	scene(api)

	for _, s := range api.surfaces {
		fmt.Println(s)
	}

	pxHigh := pxWide * api.aspectRatio[1] / api.aspectRatio[0]
	return image.NewNRGBA(image.Rect(0, 0, pxWide, pxHigh))
}
