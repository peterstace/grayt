package grayt

import "math"

type Light struct {
	Location  Vect
	Intensity float64
}

type Scene struct {
	Camera     Camera
	Geometries []Geometry
	Lights     []Light
}

func RayTracer(s Scene, a Accumulator) {

	for pxX := 0; pxX < a.wide; pxX++ {
		for pxY := 0; pxY < a.high; pxY++ {

			pxPitch := 2.0 / float64(a.wide)
			x := (float64(pxX-a.wide/2) + 0.5) * pxPitch
			y := (float64(pxY-a.high/2) + 0.5) * pxPitch * -1.0

			r := s.Camera.MakeRay(x, y)
			r.Dir = r.Dir.Unit()
			a.add(pxX, pxY, traceRay(s, r))
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
	hr, ok := closestHit(s.Geometries, r)
	if !ok {
		// Missed everything, shade black.
		return Colour{0, 0, 0}
	}

	// Subtract a small amount to the hit distance, to prevent the object
	// shaddowing itself.
	hr.Distance = addULPs(hr.Distance, -50)
	hitLoc := r.At(hr.Distance)

	var colour Colour

	// Calculate the colour at the hit point.
	for _, light := range s.Lights {

		// Vector from hit location to the light.
		fromHitToLight := light.Location.Sub(hitLoc)
		unitFromHitToLight := fromHitToLight.Unit()
		attenuation := fromHitToLight.Length2()

		// Test if anything obscures the light.
		if tmpHR, ok := closestHit(
			s.Geometries,
			Ray{Start: hitLoc, Dir: fromHitToLight},
		); !ok || tmpHR.Distance > 1.0 {

			// Lambert shading.
			lambertCoef := math.Abs(unitFromHitToLight.Dot(hr.UnitNormal))
			lambertColour := hr.Material.Colour.Scale(
				lambertCoef * light.Intensity / attenuation,
			)

			colour = colour.Add(lambertColour)
		}
	}
	return colour
}

func closestHit(gs []Geometry, r Ray) (Intersection, bool) {
	isHit := false
	var closest Intersection
	for _, geometry := range gs {
		tmpHR, ok := geometry.Intersect(r)
		if ok && (!isHit || tmpHR.Distance < closest.Distance) {
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
