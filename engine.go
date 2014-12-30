package grayt

import (
	"image"
	"image/color"
)

// TraceScene traces a single scene, given quality settings and scene samples
// (at least one sample must be provided).
func TraceScene(quality *Quality, samples ...Scene) image.Image {
	img := image.NewRGBA64(image.Rect(0, 0, quality.PxWide, quality.PxHigh))
	for pxY := 0; pxY < quality.PxHigh; pxY++ {
		for pxX := 0; pxX < quality.PxWide; pxX++ {
			// Calculate x and y
			// Calculate the ray
			// Trace the ray
			// Set the image value
			img.Set(pxX, pxY, trace())
		}
	}
	return img
}

func trace() color.Color {
	return NewColor(1, 0.5, 0)
}
