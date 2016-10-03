package main

import "github.com/peterstace/grayt/scene"

type intersection struct {
	unitNormal vector
	distance   float64

	colour    colour
	emittance float64
}

type triangle struct {
	a, u, v             vector // Corner A, A to B, and A to C.
	unitNorm            vector
	dotUV, dotUU, dotVV float64 // Precomputed dot products.

	colour    colour
	emittance float64
}

func newTriangle(a, b, c vector, colour colour, emittance float64) triangle {
	u := b.sub(a)
	v := c.sub(a)
	return triangle{
		a:        a,
		u:        u,
		v:        v,
		unitNorm: u.cross(v).unit(),
		dotUV:    u.dot(v),
		dotUU:    u.dot(u),
		dotVV:    v.dot(v),

		colour:    colour,
		emittance: emittance,
	}
}

func convertTriangles(ts []scene.Triangle) []triangle {
	tris := make([]triangle, 0, len(ts))
	for _, t := range ts {
		tris = append(tris, newTriangle(
			convertVector(t.A),
			convertVector(t.B),
			convertVector(t.C),
			convertColour(t.Colour),
			t.Emittance,
		))
	}
	return tris
}

func (t *triangle) intersect(r ray) (intersection, bool) {

	// Check if there's a hit with the plane.
	h := t.unitNorm.dot(t.a.sub(r.start)) / t.unitNorm.dot(r.dir)
	if h <= 0 {
		// Hit was behind the camera.
		return intersection{}, false
	}

	// Find out if the plane hit was inside the triangle. We need to solve the
	// equation w = alpha*u + beta*v for alpha and beta (where alpha beta
	// scalars, and u and v are vectors from a to b and a to c, and w is a
	// vector from a to the hit point).
	//
	// If the sum of alpha and beta is less than 1 and both alpha and beta are
	// positive, then the hit is inside the triangle.
	//
	// alpha = [(u.v)(w.v) - (v.v)(w.u)] / [(u.v)^2 - (u.u)(v.v)]
	// beta  = [(u.v)(w.u) - (u.u)(w.v)] / [(u.v)^2 - (u.u)(v.v)]

	w := r.at(h).sub(t.a)
	dotWV := w.dot(t.v)
	dotWU := w.dot(t.u)
	alpha := t.dotUV*dotWV - t.dotVV*dotWU
	beta := t.dotUV*dotWU - t.dotUU*dotWV
	denom := t.dotUV*t.dotUV - t.dotUU*t.dotVV
	alpha /= denom
	beta /= denom

	if alpha < 0 || beta < 0 || alpha+beta > 1 {
		return intersection{}, false
	}
	return intersection{
		unitNormal: t.unitNorm,
		distance:   h,
		colour:     t.colour,
		emittance:  t.emittance,
	}, true
}
