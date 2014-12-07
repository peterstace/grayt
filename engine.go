package grayt

import (
	"image"
	"image/color"
)

type engine struct {
	quality Quality
}

func (e *engine) traceScenes(s []Scene) image.Image {
	img := image.NewRGBA64(image.Rect(0, 0, e.quality.PxWide, e.quality.PxHigh))
	for pxY := 0; pxY < e.quality.PxHigh; pxY++ {
		for pxX := 0; pxX < e.quality.PxWide; pxX++ {
			// Calculate x and y
			// Calculate the ray
			// Trace the ray
			// Set the image value
			img.Set(pxX, pxY, e.trace())
		}
	}
	return img
}

func (e *engine) trace() color.Color {
	return NewColor(1, 0.5, 0)
}
