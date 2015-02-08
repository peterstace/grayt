package tracer

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/peterstace/grayt/ray"
)

func TraceImage(samples []Scene) image.Image {

	const pxWide = 640
	const pxHigh = 480

	img := image.NewGray(image.Rect(0, 0, pxWide, pxHigh))

	for pxX := 0; pxX < pxWide; pxX++ {
		for pxY := 0; pxY < pxHigh; pxY++ {

			pxPitch := 2.0 / float64(pxWide)
			x := (float64(pxX-pxWide/2) + 0.5) * pxPitch
			y := (float64(pxY-pxHigh/2) + 0.5) * pxPitch

			s := &samples[rand.Intn(len(samples))]
			r := s.Camera.MakeRay(x, y)
			img.Set(pxX, pxY, traceRay(s, r))
		}
	}

	return img
}

func traceRay(s *Scene, r ray.Ray) color.Color {
	return nil
}
