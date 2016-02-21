package geometry

import "github.com/peterstace/grayt"

func Triangle(a, b, c grayt.Vector) grayt.Surface {
	u := b.Sub(a)
	v := c.Sub(a)
	return &triangle{
		a:        a,
		u:        u,
		v:        v,
		unitNorm: u.Cross(v).Unit(),
		dotUV:    u.Dot(v),
		dotUU:    u.Dot(u),
		dotVV:    v.Dot(v),
	}
}

type triangle struct {
	a, u, v             grayt.Vector // Corner A, A to B, and A to C.
	unitNorm            grayt.Vector
	dotUV, dotUU, dotVV float64 // Precomputed dot products.
}

func (t *triangle) Intersect(r grayt.Ray) (grayt.Intersection, bool) {

	// Check if there's a hit with the plane.
	h := t.unitNorm.Dot(t.a.Sub(r.Start)) / t.unitNorm.Dot(r.Dir)
	if h <= 0 {
		// Hit was behind the camera.
		return grayt.Intersection{}, false
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

	w := r.At(h).Sub(t.a)
	dotWV := w.Dot(t.v)
	dotWU := w.Dot(t.u)
	alpha := t.dotUV*dotWV - t.dotVV*dotWU
	beta := t.dotUV*dotWU - t.dotUU*dotWV
	denom := t.dotUV*t.dotUV - t.dotUU*t.dotVV
	alpha /= denom
	beta /= denom

	if alpha < 0 || beta < 0 || alpha+beta > 1 {
		return grayt.Intersection{}, false
	}
	return grayt.Intersection{UnitNormal: t.unitNorm, Distance: h}, true
}
