package grayt

import "math/rand"

func TraceImage(c Camera, w *world, acc accumulator) {
	pxPitch := 2.0 / float64(acc.wide)
	for pxX := 0; pxX < acc.wide; pxX++ {
		for pxY := 0; pxY < acc.high; pxY++ {
			//if pxX != 140 || pxY != 40 {
			//	continue
			//}
			x := (float64(pxX-acc.wide/2) + rand.Float64()) * pxPitch
			y := (float64(pxY-acc.high/2) + rand.Float64()) * pxPitch * -1.0
			r := c.MakeRay(x, y)
			r.Dir = r.Dir.Unit()
			acc.add(pxX, pxY, tracePath(w, r))
		}
	}
}

func tracePath(w *world, r Ray) Colour {

	intersection, material := w.closestHit(r)
	if material == nil {
		return Colour{0, 0, 0}
	}

	// Calculate probability of emitting.
	pEmit := 0.1
	if material.Emittance != 0 {
		pEmit = 1.0
	}

	// Handle emit case.
	if rand.Float64() < pEmit {
		return material.Colour.
			Scale(material.Emittance / pEmit)
	}

	// Find where the ray hit. Reduce the intersection distance by a small
	// amount so that reflected rays don't intersect with it immediately.
	hitLoc := r.At(addULPs(intersection.Distance, -50))

	// Orient the unit normal towards the ray origin.
	if intersection.UnitNormal.Dot(r.Dir) > 0 {
		intersection.UnitNormal = intersection.UnitNormal.Scale(-1.0)
	}

	// Create a random vector on the hemisphere towards the normal.
	rnd := Vector{rand.NormFloat64(), rand.NormFloat64(), rand.NormFloat64()}
	rnd = rnd.Unit()
	if rnd.Dot(intersection.UnitNormal) < 0 {
		rnd = rnd.Scale(-1.0)
	}

	// Apply the BRDF (bidirectional reflection distribution function).
	brdf := rnd.Dot(intersection.UnitNormal)

	return tracePath(w, Ray{Start: hitLoc, Dir: rnd}).
		Scale(brdf / (1 - pEmit)).
		Mul(material.Colour)
}
