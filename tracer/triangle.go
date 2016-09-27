package main

import "github.com/peterstace/grayt/scene"

type intersection struct {
	unitNormal vector
	distance   float64
}

type triangle struct {
	a, u, v             vector // Corner A, A to B, and A to C.
	unitNorm            vector
	dotUV, dotUU, dotVV float64 // Precomputed dot products.
}

func newTriangle(a, b, c vector) triangle {
	return triangle{
		a:        a,
		u:        u,
		v:        v,
		unitNorm: u.Cross(v).Unit(),
		dotUV:    u.Dot(v),
		dotUU:    u.Dot(u),
		dotVV:    v.Dot(v),
	}
}

func convertTriangles(ts []scene.Triangle) []triangle {
	tris := make([]triangle, 0, len(s.Triangles))
	for _, t := range s.Triangles {
		cornerA := vector{t.A.X, t.A.Y, t.A.Z}
		cornerB := vector{t.B.X, t.B.Y, t.B.Z}
		cornerC := vector{t.C.X, t.C.Y, t.C.Z}
		tris = append(tris, newTriangle(cornerA, cornerB, cornerC))
	}
	return tris
}

func (t *triangle) intersect(r ray) (intersection, bool) {

	// Check if there's a hit with the plane.
	h := t.unitNorm.Dot(t.a.Sub(r.Start)) / t.unitNorm.Dot(r.Dir)
	if h <= 0 {
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
	return Intersection{UnitNormal: t.unitNorm, Distance: h}, true
}
