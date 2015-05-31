package grayt

import (
	"math"
)

func RayTracer(s Scene, a Accumulator, samplesPerPixel int) {

	for pxX := 0; pxX < a.wide; pxX++ {
		for pxY := 0; pxY < a.high; pxY++ {

			//if pxX != 140 || pxY != 40 {
			//	continue
			//}

			pxPitch := 2.0 / float64(a.wide)
			x := (float64(pxX-a.wide/2) + 0.5) * pxPitch
			y := (float64(pxY-a.high/2) + 0.5) * pxPitch * -1.0

			for i := 0; i < samplesPerPixel; i++ {
				r := s.Camera.MakeRay(x, y)
				r.Dir = r.Dir.Unit()
				a.add(pxX, pxY, traceRay(s, r))
			}
		}
	}
}

// traceRay is a recursive function to find the colour from a single ray into a
// scene. It's a precondition that r.Dir must be a unit vector.
func traceRay(s Scene, r Ray) Colour {

	// Assert that r.Dir is a unit vector.
	if ulpDiff(1.0, r.Dir.Length2()) > 50 {
		panic("precondition not met: r.Dir not a unit vector")
	}

	// Establish the hit point.
	intersection, reflector := closestHit(s.Reflectors, r)
	if reflector == nil {
		// Missed everything, shade black.
		return Colour{0, 0, 0}
	}

	// Subtract a small amount to the hit distance, to prevent the object
	// shaddowing itself.
	intersection.Distance = addULPs(intersection.Distance, -50)
	hitLoc := r.At(intersection.Distance)

	var colour Colour

	// Calculate the colour at the hit point.
	for _, emitter := range s.Emitters {

		// Vector from hit location to the light.
		fromHitToLight := emitter.Sample().Sub(hitLoc)
		unitFromHitToLight := fromHitToLight.Unit()
		attenuation := fromHitToLight.Length2()

		// Test if anything obscures the light.
		if maskIntersection, mask := closestHit(
			s.Reflectors,
			Ray{Start: hitLoc, Dir: fromHitToLight},
		); mask == nil || maskIntersection.Distance > 1.0 {

			// Lambert shading.
			lambertCoef := math.Abs(unitFromHitToLight.Dot(intersection.UnitNormal))
			lambertColour := reflector.Material.Colour.Scale(
				lambertCoef * emitter.Intensity / attenuation,
			)
			colour = colour.Add(lambertColour)
		}
	}

	// Add ambient light.
	const ambientCoef = 0.3
	colour = colour.Scale(1 - ambientCoef).Add(reflector.Material.Colour.Scale(ambientCoef))

	return colour
}

func closestHit(reflectors []Reflector, r Ray) (Intersection, *Reflector) {
	var closest struct {
		Intersection
		*Reflector
	}
	for i := range reflectors {
		intersection, hit := reflectors[i].Intersect(r)
		if !hit {
			continue
		}
		if closest.Reflector == nil || intersection.Distance < closest.Intersection.Distance {
			closest.Intersection = intersection
			closest.Reflector = &reflectors[i]
		}
	}
	return closest.Intersection, closest.Reflector
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
