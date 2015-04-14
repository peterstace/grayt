package tracer

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/peterstace/grayt/ray"
	"github.com/peterstace/grayt/vect"
)

func TraceImage(samples []Scene) image.Image {

	const pxWide = 640
	const pxHigh = 480

	startTime := time.Now()
	defer func() {
		totalTime := time.Since(startTime)
		pxCount := pxWide * pxHigh
		timePerPixel := time.Duration(int(totalTime) / pxCount)

		log.Printf("Dimensions=%dx%d PxCount=%d TotalTime=%s TimePerPixel=%s",
			pxWide, pxHigh, pxCount, totalTime, timePerPixel)
	}()

	img := image.NewGray(image.Rect(0, 0, pxWide, pxHigh))

	for pxX := 0; pxX < pxWide; pxX++ {
		for pxY := 0; pxY < pxHigh; pxY++ {

			pxPitch := 2.0 / float64(pxWide)
			x := (float64(pxX-pxWide/2) + 0.5) * pxPitch
			y := (float64(pxY-pxHigh/2) + 0.5) * pxPitch * -1.0

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
	hr, ok := closestHit(s.Geometries, r)
	if !ok {
		// Missed everything, shade black.
		return color.Gray{Y: 0x00}
	}

	hitLoc := r.At(hr.distance * 0.999999) // XXX move by several ULPs

	var colour float64

	// Calculate the colour at the hit point.
	for _, light := range s.Lights {

		// Vector from hit location to the light.
		fromHitToLight := vect.Sub(light.sampleLocation(), hitLoc)
		unitFromHitToLight := fromHitToLight.Unit()
		attenuation := fromHitToLight.Length2()

		// Test if anything obscures the light.
		if tmpHR, ok := closestHit(
			s.Geometries,
			ray.Ray{Start: hitLoc, Dir: fromHitToLight},
		); !ok || tmpHR.distance > 1.0 {

			// Lambert shading.
			lambertCoef := vect.Dot(unitFromHitToLight, hr.unitNormal)

			// Add shading to the colour.
			colour += math.Abs(lambertCoef) * light.Intensity / attenuation
		}
	}

	var colourUint8 uint8
	if colour >= 1.0 {
		colourUint8 = 0xff
	} else if colour < 0.0 {
		colourUint8 = 0x00
	} else {
		colourUint8 = uint8(colour * 256)
	}
	return color.Gray{Y: colourUint8}
}

func closestHit(gs []Geometry, r ray.Ray) (hitRec, bool) {
	isHit := false
	var closest hitRec
	for _, geometry := range gs {
		tmpHR, ok := geometry.intersect(r)
		if ok && (!isHit || tmpHR.distance < closest.distance) {
			closest = tmpHR
			isHit = true
		}
	}
	return closest, isHit
}
