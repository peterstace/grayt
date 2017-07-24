package grayt

import (
	"fmt"
	"math"
)

type surface interface {
	intersect(r ray) (intersection, bool)
	bound() (Vector, Vector)
	translate(Vector)
	rotate(Vector, float64)
	scale(float64)
}

type material struct {
	colour    Colour
	emittance float64
	mirror    bool
}

type Object struct {
	surface
	material
}

func (o Object) String() string {
	return fmt.Sprintf("Surface={%v} Material={%v}", o.surface, o.material)
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

func (t *triangle) String() string {
	return fmt.Sprintf("Type=triangle A=%v B=%v C=%v", t.a, t.a.Add(t.u), t.a.Add(t.v))
}

func newTriangle(a, b, c Vector) *triangle {
	u := b.Sub(a)
	v := c.Sub(a)
	return &triangle{
		a:        a,
		u:        u,
		v:        v,
		unitNorm: u.cross(v).Unit(),
		dotUV:    u.Dot(v),
		dotUU:    u.Dot(u),
		dotVV:    v.Dot(v),
	}
}

func (t *triangle) intersect(r ray) (intersection, bool) {

	// Check if there's a hit with the plane.
	h := t.unitNorm.Dot(t.a.Sub(r.start)) / t.unitNorm.Dot(r.dir)
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
	dotWV := w.Dot(t.v)
	dotWU := w.Dot(t.u)
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

func (t *triangle) bound() (Vector, Vector) {
	b := t.a.Add(t.u)
	c := t.a.Add(t.v)
	min := t.a.Min(b.Min(c)).addULPs(-ulpFudgeFactor)
	max := t.a.Max(b.Max(c)).addULPs(+ulpFudgeFactor)
	return min, max
}

func (t *triangle) translate(v Vector) {
	a := t.a.Add(v)
	b := t.u.Add(t.a).Add(v)
	c := t.v.Add(t.a).Add(v)
	*t = *newTriangle(a, b, c)
}

func (t *triangle) rotate(v Vector, rads float64) {
	a := t.a.rotate(v, rads)
	b := t.u.Add(t.a).rotate(v, rads)
	c := t.v.Add(t.a).rotate(v, rads)
	*t = *newTriangle(a, b, c)
}

func (t *triangle) scale(f float64) {
	a := t.a.Scale(f)
	b := t.u.Add(t.a).Scale(f)
	c := t.v.Add(t.a).Scale(f)
	*t = *newTriangle(a, b, c)
}

type alignedBox struct {
	min, max Vector
}

func (a *alignedBox) String() string {
	return fmt.Sprintf("Type=alignedBox Min=%v Max=%v", a.min, a.max)
}

func newAlignedBox(corner1, corner2 Vector) surface {
	return &alignedBox{
		min: corner1.Min(corner2),
		max: corner1.Max(corner2),
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
		if ty1 < ty2 {
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
		if tz1 < tz2 {
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

func (b *alignedBox) bound() (Vector, Vector) {
	return b.min, b.max
}

func (b *alignedBox) translate(v Vector) {
	b.min = b.min.Add(v)
	b.max = b.max.Add(v)
}

func (b *alignedBox) rotate(Vector, float64) {
	panic("cannot rotate aligned box")
}

func (b *alignedBox) scale(f float64) {
	b.min = b.min.Scale(f)
	b.max = b.max.Scale(f)
}

type sphere struct {
	centre Vector
	radius float64
}

func (s *sphere) String() string {
	return fmt.Sprintf("Type=sphere C=%v R=%v", s.centre, s.radius)
}

func (s *sphere) intersect(r ray) (intersection, bool) {

	// Get coeficients to a.x^2 + b.x + c = 0
	emc := r.start.Sub(s.centre)
	a := r.dir.LengthSq()
	b := 2 * emc.Dot(r.dir)
	c := emc.LengthSq() - s.radius*s.radius

	// Find discrimenant b*b - 4*a*c
	disc := b*b - 4*a*c
	if disc < 0 {
		return intersection{}, false
	}

	// Find x1 and x2 using a numerically stable algorithm.
	var signOfB float64
	signOfB = math.Copysign(1.0, b)
	q := -0.5 * (b + signOfB*math.Sqrt(disc))
	x1 := q / a
	x2 := c / q

	var t float64
	if x1 > 0 && x2 > 0 {
		// Both are positive, so take the smaller one.
		t = math.Min(x1, x2)
	} else {
		// At least one is negative, take the larger one (which is either
		// negative or positive).
		t = math.Max(x1, x2)
	}

	return intersection{
		unitNormal: r.at(t).Sub(s.centre).Unit(),
		distance:   t,
	}, t > 0
}

func (s *sphere) bound() (Vector, Vector) {
	r := Vect(s.radius, s.radius, s.radius)
	min, max := s.centre.Sub(r), s.centre.Add(r)
	return min.addULPs(-ulpFudgeFactor), max.addULPs(ulpFudgeFactor)
}

func (s *sphere) translate(v Vector) {
	s.centre = s.centre.Add(v)
}

func (s *sphere) rotate(v Vector, rads float64) {
	// NO-OP
}

func (s *sphere) scale(f float64) {
	s.radius *= f
	s.centre = s.centre.Scale(f)
}

type alignXSquare struct {
	x, y1, y2, z1, z2 float64
}

func (a *alignXSquare) String() string {
	return fmt.Sprintf("Type=alignXSquare X=%v Y1=%v Y2=%v Z1=%v Z2=%v",
		a.x, a.y1, a.y2, a.z1, a.z2)
}

func (s *alignXSquare) intersect(r ray) (intersection, bool) {
	t := (s.x - r.start.X) / r.dir.X
	hit := r.at(t)
	return intersection{Vect(+1, 0, 0), t},
		t > 0 && hit.Y > s.y1 && hit.Y < s.y2 && hit.Z > s.z1 && hit.Z < s.z2
}

func (s *alignXSquare) bound() (Vector, Vector) {
	return Vect(s.x, s.y1, s.z1), Vect(s.x, s.y2, s.z2)
}

func (s *alignXSquare) translate(v Vector) {
	s.x += v.X
	s.y1 += v.Y
	s.y2 += v.Y
	s.z1 += v.Z
	s.z2 += v.Z
}

func (a *alignXSquare) rotate(Vector, float64) {
	panic("cannot rotate aligned square")
}

func (a *alignXSquare) scale(f float64) {
	a.x *= f
	a.y1 *= f
	a.y2 *= f
	a.z1 *= f
	a.z2 *= f
}

type alignYSquare struct {
	x1, x2, y, z1, z2 float64
}

func (a *alignYSquare) String() string {
	return fmt.Sprintf("Type=alignYSquare X1=%v X2=%v Y=%v Z1=%v Z2=%v",
		a.x1, a.x2, a.y, a.z1, a.z2)
}

func (s *alignYSquare) intersect(r ray) (intersection, bool) {
	t := (s.y - r.start.Y) / r.dir.Y
	hit := r.at(t)
	return intersection{Vect(0, +1, 0), t},
		t > 0 && hit.X > s.x1 && hit.X < s.x2 && hit.Z > s.z1 && hit.Z < s.z2
}

func (s *alignYSquare) bound() (Vector, Vector) {
	return Vect(s.x1, s.y, s.z1), Vect(s.x2, s.y, s.z2)
}

func (s *alignYSquare) translate(v Vector) {
	s.x1 += v.X
	s.x2 += v.X
	s.y += v.Y
	s.z1 += v.Z
	s.z2 += v.Z
}

func (a *alignYSquare) rotate(Vector, float64) {
	panic("cannot rotate aligned square")
}

func (a *alignYSquare) scale(f float64) {
	a.x1 *= f
	a.x2 *= f
	a.y *= f
	a.z1 *= f
	a.z2 *= f
}

type alignZSquare struct {
	x1, x2, y1, y2, z float64
}

func (a *alignZSquare) String() string {
	return fmt.Sprintf("Type=alignZSquare X1=%v X2=%v Y1=%v Y2=%v Z=%v",
		a.x1, a.x2, a.y1, a.y2, a.z)
}

func (s *alignZSquare) intersect(r ray) (intersection, bool) {
	t := (s.z - r.start.Z) / r.dir.Z
	hit := r.at(t)
	return intersection{Vect(0, 0, +1), t},
		t > 0 && hit.X > s.x1 && hit.X < s.x2 && hit.Y > s.y1 && hit.Y < s.y2
}

func (s *alignZSquare) bound() (Vector, Vector) {
	return Vect(s.x1, s.y1, s.z), Vect(s.x2, s.y2, s.z)
}

func (s *alignZSquare) translate(v Vector) {
	s.x1 += v.X
	s.x2 += v.X
	s.y1 += v.Y
	s.y2 += v.Y
	s.z += v.Z
}

func (a *alignZSquare) rotate(Vector, float64) {
	panic("cannot rotate aligned square")
}

func (a *alignZSquare) scale(f float64) {
	a.x1 *= f
	a.x2 *= f
	a.y1 *= f
	a.y2 *= f
	a.z *= f
}
