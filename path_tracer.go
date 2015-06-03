package grayt

import "math/rand"

// PathTracer traces a scene using the path tracing algorithm.
func PathTracer(s Scene, a Accumulator, spp int) {
	surfs := make([]Surface, len(s.Emitters)+len(s.Reflectors))
	for i := range s.Reflectors {
		surfs[i] = s.Reflectors[i]
	}
	for i := range s.Emitters {
		surfs[len(s.Reflectors)+i] = s.Emitters[i]
	}
	trace(func(r Ray) Colour {
		return tracePath(surfs, r, 0)
	}, s.Camera, a, spp)
}

func tracePath(surfs []Surface, r Ray, i int) Colour {

	if i == 10 {
		return Colour{0, 0, 0}
	}

	intersection, hitSurf := closestHit(surfs, r)
	if hitSurf == nil {
		return Colour{0, 0, 0}
	}

	if emitter, ok := hitSurf.(Emitter); ok {
		return emitter.Colour.Scale(emitter.Intensity)
	}
	reflector := hitSurf.(Reflector)

	rnd := Vect{rand.NormFloat64(), rand.NormFloat64(), rand.NormFloat64()}
	rnd = rnd.Unit()
	if rnd.Dot(intersection.UnitNormal) < 0 {
		rnd.Extended(-1)
	}
	hitLoc := r.At(addULPs(intersection.Distance, -50))
	colour := tracePath(surfs, Ray{Start: hitLoc, Dir: rnd}, i+1)
	colour = colour.Mul(reflector.Material.Colour)
	return colour
}
