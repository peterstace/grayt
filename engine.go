package grayt

import (
	"image"
	"image/color"
	"math"
)

type engine struct {
	quality Quality
}

func (e *engine) traceScenes(s []Scene) image.Image {
	img := image.NewRGBA64(image.Rect(0, 0, e.quality.PxWide, e.quality.PxHigh))
	for y := 0; y < e.quality.PxHigh; y++ {
		for x := 0; x < e.quality.PxWide; x++ {
			img.Set(x, y, e.trace())
		}
	}
	return img
}

func (e *engine) trace() color.Color {
	c := color.NRGBA64{
		R: math.MaxUint16,
		G: math.MaxUint16 / 2,
		B: 0,
		A: math.MaxUint16,
	}
	return c
}
