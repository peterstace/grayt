package grayt

type surface interface {
	intersect(r ray) (intersection, bool)
}

type material struct {
	colour    Colour
	emittance float64
}

type Object struct {
	surface
	material
}

type intersection struct {
	unitNormal Vector
	distance   float64
}

type triangle struct {
	a, u, v             Vector // Corner A, A to B, and A to C.
	unitNorm            Vector
	dotUV, dotUU, dotVV float64 // Precomputed dot products.
}

func newTriangle(a, b, c Vector) surface {
	u := b.Sub(a)
	v := c.Sub(a)
	return &triangle{
		a:        a,
		u:        u,
		v:        v,
		unitNorm: u.cross(v).unit(),
		dotUV:    u.dot(v),
		dotUU:    u.dot(u),
		dotVV:    v.dot(v),
	}
}

func (t *triangle) intersect(r ray) (intersection, bool) {

	// Check if there's a hit with the plane.
	h := t.unitNorm.dot(t.a.Sub(r.start)) / t.unitNorm.dot(r.dir)
	if h <= 0 {
		// Hit was behind the camera.
		return intersection{}, false
	}

	// Find out if the plane hit was inside the triangle. We need to solve the
	// equation w = alpha*u + beta*v for alpha and beta (where alpha and beta
	// are scalars, and u and v are vectors from a to b and a to c, and w is a
	// vector from a to the hit point).
	//
	// If the sum of alpha and beta is less than 1 and both alpha and beta are
	// positive, then the hit is inside the triangle.
	//
	// alpha = [(u.v)(w.v) - (v.v)(w.u)] / [(u.v)^2 - (u.u)(v.v)]
	// beta  = [(u.v)(w.u) - (u.u)(w.v)] / [(u.v)^2 - (u.u)(v.v)]

	w := r.at(h).Sub(t.a)
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
	}, true
}
