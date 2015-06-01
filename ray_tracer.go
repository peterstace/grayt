package grayt

import "math"

type rayTracerWorld struct {
	reflectors []Surface // always stores reflectors.
	emitters   []Emitter
}

// RayTracer traces a scene using the distributed ray tracing algorith. The
// ambient shading colour is white with an intensity of 1.0.
func RayTracer(s Scene, a Accumulator, samplesPerPixel int) {
	world := rayTracerWorld{
		reflectors: make([]Surface, len(s.Reflectors)),
		emitters:   s.Emitters,
	}
	for i := range s.Reflectors {
		world.reflectors[i] = s.Reflectors[i]
	}
	trace(func(r Ray) Colour {
		return traceRay(world, r)
	}, s.Camera, a, samplesPerPixel)
}

// traceRay is a recursive function to find the colour from a single ray into a
// scene. It's a precondition that r.Dir must be a unit vector.
func traceRay(w rayTracerWorld, r Ray) Colour {

	// Assert that r.Dir is a unit vector.
	if ulpDiff(1.0, r.Dir.Length2()) > 50 {
		panic("precondition not met: r.Dir not a unit vector")
	}

	// Establish the hit point.
	intersection, reflectorSurface := closestHit(w.reflectors, r)
	if reflectorSurface == nil {
		// Missed everything, shade black.
		return Colour{0, 0, 0}
	}
	reflector := reflectorSurface.(Reflector)

	// Subtract a small amount to the hit distance, to prevent the object
	// shaddowing itself.
	intersection.Distance = addULPs(intersection.Distance, -50)
	hitLoc := r.At(intersection.Distance)

	var colour Colour

	// Calculate the colour at the hit point.
	for _, emitter := range w.emitters {

		// Vector from hit location to the light.
		fromHitToLight := emitter.Sample().Sub(hitLoc)
		unitFromHitToLight := fromHitToLight.Unit()
		attenuation := fromHitToLight.Length2()

		// Test if anything obscures the light.
		if maskIntersection, mask := closestHit(
			w.reflectors,
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

	// Add ambient light (white, intensity of 1.0).
	colour = colour.Add(reflector.Material.Colour)

	return colour
}
