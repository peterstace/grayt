package grayt

import "image"

type engine struct {
	config Config
}

func (e *engine) traceScenes(s []Scene) image.Image {
	return image.NewRGBA64(image.Rect(0, 0, e.config.PxWide, e.config.PxHigh))
}
