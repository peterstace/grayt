package trace

import (
	"math"
	"math/rand"

	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/xmath"
)

func newTracer(accel accelerationStructure, rng *rand.Rand) *tracer {
	return &tracer{accel: accel, rng: rng}
}

type tracer struct {
	accel accelerationStructure
	rng   *rand.Rand
}

func (t *tracer) tracePath(r xmath.Ray) colour.Colour {
	assertUnit(r.Dir)
	intersection, material, hit := t.accel.closestHit(r)
	if !hit {
		return colour.Colour{0, 0, 0}
	}
	assertUnit(intersection.unitNormal)

	// Calculate probability of emitting.
	pEmit := 0.1
	if material.Emittance != 0 {
		pEmit = 1.0
	}

	// Handle emit case.
	if t.rng.Float64() < pEmit {
		return material.Colour.Scale(material.Emittance / pEmit)
	}

	offsetScale := -math.Copysign(xmath.AddULPs(1.0, 1e5)-1.0, r.Dir.Dot(intersection.unitNormal))
	offset := intersection.unitNormal.Scale(offsetScale)
	hitLoc := r.At(intersection.distance).Add(offset)

	// Orient the unit normal towards the ray origin.
	if intersection.unitNormal.Dot(r.Dir) > 0 {
		intersection.unitNormal = intersection.unitNormal.Scale(-1.0)
	}

	if material.Mirror {

		reflected := r.Dir.Sub(intersection.unitNormal.Scale(2 * intersection.unitNormal.Dot(r.Dir)))
		return t.tracePath(xmath.Ray{Start: hitLoc, Dir: reflected})

	} else {

		// Create a random vector on the hemisphere towards the normal.
		rnd := xmath.Vector{t.rng.NormFloat64(), t.rng.NormFloat64(), t.rng.NormFloat64()}
		rnd = rnd.Unit()
		if rnd.Dot(intersection.unitNormal) < 0 {
			rnd = rnd.Scale(-1.0)
		}

		// Apply the BRDF (bidirectional reflection distribution function).
		brdf := rnd.Dot(intersection.unitNormal)

		return t.tracePath(xmath.Ray{Start: hitLoc, Dir: rnd}).
			Scale(brdf / (1 - pEmit)).
			Mul(material.Colour)
	}
}
