package grayt

func NewTriangle(cornerA, cornerB, cornerC Vect) Surface {
	u := cornerB.Sub(cornerA)
	v := cornerC.Sub(cornerA)
	t := &triangle{
		a:        cornerA,
		u:        u,
		v:        v,
		unitNorm: u.Cross(v).Unit(),
		dotUV:    u.Dot(v),
		dotUU:    u.Dot(u),
		dotVV:    v.Dot(v),
	}
	return t
}

type triangle struct {
	a, u, v             Vect // Corner A, A to B, and A to C.
	unitNorm            Vect
	dotUV, dotUU, dotVV float64 // Precomputed dot products.
}

func (t *triangle) Intersect(r Ray) (Intersection, bool) {

	// Check if there's a hit with the plane.
	h := t.unitNorm.Dot(t.a.Sub(r.Start)) / t.unitNorm.Dot(r.Dir)
	if h < 0 {
		// Hit was behind the camera.
		return Intersection{}, false
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
		return Intersection{}, false
	}
	return Intersection{t.unitNorm, h}, true
}
