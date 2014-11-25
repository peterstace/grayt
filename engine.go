package grayt

import "image"

type Engine struct {
	pxWide, pxHigh int
}

func newEngine() Engine {
	return Engine{
		pxWide: 640,
		pxHigh: 480,
	}
}

func (e *Engine) traceScenes(s []Scene) image.Image {
	return image.NewRGBA64(image.Rect(0, 0, e.pxWide, e.pxHigh))
}
