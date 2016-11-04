package grayt

import "math"

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

type alignedBox struct {
	min, max Vector
}

func newAlignedBox(corner1, corner2 Vector) surface {
	return &alignedBox{
		min: corner1.min(corner2),
		max: corner1.max(corner2),
	}
}

func (b *alignedBox) intersect(r ray) (intersection, bool) {

	tx1 := (b.min.X - r.start.X) / r.dir.X
	tx2 := (b.max.X - r.start.X) / r.dir.X
	ty1 := (b.min.Y - r.start.Y) / r.dir.Y
	ty2 := (b.max.Y - r.start.Y) / r.dir.Y
	tz1 := (b.min.Z - r.start.Z) / r.dir.Z
	tz2 := (b.max.Z - r.start.Z) / r.dir.Z

	tmin, tmax := math.Inf(-1), math.Inf(+1)
	var nMin Vector
	var nMax Vector

	if math.Min(tx1, tx2) > tmin {
		if tx1 < tx2 {
			tmin = tx1
			nMin = Vect(-1, 0, 0)
		} else {
			tmin = tx2
			nMin = Vect(1, 0, 0)
		}
	}
	if math.Max(tx1, tx2) < tmax {
		if tx1 > tx2 {
			tmax = tx1
			nMax = Vect(-1, 0, 0)
		} else {
			tmax = tx2
			nMax = Vect(1, 0, 0)
		}
	}

	if math.Min(ty1, ty2) > tmin {
		if ty1 < ty2 && ty1 > 0 {
			tmin = ty1
			nMin = Vect(0, -1, 0)
		} else {
			tmin = ty2
			nMin = Vect(0, 1, 0)
		}
	}
	if math.Max(ty1, ty2) < tmax {
		if ty1 > ty2 {
			tmax = ty1
			nMax = Vect(0, -1, 0)
		} else {
			tmax = ty2
			nMax = Vect(0, 1, 0)
		}
	}

	if math.Min(tz1, tz2) > tmin {
		if tz1 < tz2 && tz1 > 0 {
			tmin = tz1
			nMin = Vect(0, 0, -1)
		} else {
			tmin = tz2
			nMin = Vect(0, 0, 1)
		}
	}
	if math.Max(tz1, tz2) < tmax {
		if tz1 > tz2 {
			tmax = tz1
			nMax = Vect(0, 0, -1)
		} else {
			tmax = tz2
			nMax = Vect(0, 0, 1)
		}
	}

	if tmin > tmax || tmax <= 0 {
		return intersection{}, false
	}

	if tmin > 0 {
		return intersection{distance: tmin, unitNormal: nMin}, true
	} else {
		return intersection{distance: tmax, unitNormal: nMax}, true
	}
}
