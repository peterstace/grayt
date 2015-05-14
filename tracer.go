package grayt

import (
	"image"
	"image/color"
	"log"
	"math"
	"time"
)

type Scene struct {
	Camera     *camera
	Geometries []geometry
	Lights     []Light
}

func TraceImage(s Scene) image.Image {

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

			r := s.Camera.MakeRay(x, y)
			r.Dir = r.Dir.Unit()
			img.Set(pxX, pxY, traceRay(s, r))
		}
	}

	return img
}

// traceRay is a recursive function to find the colour from a single ray into a
// scene. It's a precondition that r.Dir must be a unit vector.
func traceRay(s Scene, r Ray) color.Color {

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
		fromHitToLight := light.sampleLocation().Sub(hitLoc)
		unitFromHitToLight := fromHitToLight.Unit()
		attenuation := fromHitToLight.Length2()

		// Test if anything obscures the light.
		if tmpHR, ok := closestHit(
			s.Geometries,
			Ray{Start: hitLoc, Dir: fromHitToLight},
		); !ok || tmpHR.distance > 1.0 {

			// Lambert shading.
			lambertCoef := unitFromHitToLight.Dot(hr.unitNormal)

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

func closestHit(gs []geometry, r Ray) (intersection, bool) {
	isHit := false
	var closest intersection
	for _, geometry := range gs {
		tmpHR, ok := geometry.intersect(r)
		if ok && (!isHit || tmpHR.distance < closest.distance) {
			closest = tmpHR
			isHit = true
		}
	}
	return closest, isHit
}

func ulpDiff(a, b float64) uint64 {

	ulpA := math.Float64bits(a)
	ulpB := math.Float64bits(b)

	if ulpA > ulpB {
		return ulpA - ulpB
	} else {
		return ulpB - ulpA
	}
}
