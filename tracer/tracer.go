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

			s := samples[rand.Intn(len(samples))]
			r := s.Camera.MakeRay(x, y)
			r.Dir = r.Dir.Unit()
			img.Set(pxX, pxY, traceRay(s, r))
		}
	}

	return img
}

// traceRay is a recursive function to find the colour from a single ray into a
// scene.
//
// Preconditions:
//  * r.Dir must be a unit vector.
func traceRay(s Scene, r ray.Ray) color.Color {

	// Assert that r.Dir is a unit vector.
	if ulpDiff(1.0, r.Dir.Length2()) > 50 {
		panic("precondition not met: r.Dir not a unit vector")
	}

	// Establish the hit point.
	_, ok := closestHit(s.Geometries, r)
	if !ok {
		// Missed everything, shade black.
		return color.Gray{Y: 0x00}
	}

	// Hit something, shade white.
	return color.Gray{Y: 0xff}
}

func closestHit(gs []Geometry, r ray.Ray) (hitRec, bool) {
	isHit := false
	var closest hitRec
	for _, geometry := range gs {
		tmpHR, ok := geometry.intersect(r)
		if ok && (!isHit || tmpHR.t < closest.t) {
			closest = tmpHR
			isHit = true
		}
	}
	return closest, isHit
}
