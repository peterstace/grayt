package grayt

import "image"

type Engine struct {
	pxWide, pxHigh int
}

func (e *Engine) SetImageDimensions(pxWide, pxHigh int) {
	e.pxWide = pxWide
	e.pxHigh = pxHigh
}

func (e *Engine) traceScenes(s []Scene) image.Image {
	return image.NewRGBA64(image.Rect(0, 0, e.pxWide, e.pxHigh))
}
