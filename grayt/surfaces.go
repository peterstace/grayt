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
	Colour    Colour  `json:"colour"`
	Emittance float64 `json:"emittance"`
	Mirror    bool    `json:"mirror"`
}

type Object struct {
	Surface  surface  `json:"surface"`
	Material material `json:"material"`
}

func (o Object) String() string {
	return fmt.Sprintf("Surface={%v} Material={%v}", o.Surface, o.Material)
}

type intersection struct {
	unitNormal Vector
	distance   float64
}

type triangle struct {
	A        Vector `json:"a"` // Corner A
	U        Vector `json:"u"` // A to B
	V        Vector `json:"v"` // A to C
	UnitNorm Vector `json:"unit_norm"`
	// Precomputed dot products:
	DotUV float64 `json:"dot_uv"`
	DotUU float64 `json:"dot_uu"`
	DotVV float64 `json:"dot_vv"`
}

func (t *triangle) String() string {
	return fmt.Sprintf("Type=triangle A=%v B=%v C=%v", t.A, t.A.Add(t.U), t.A.Add(t.V))
}

func newTriangle(a, b, c Vector) *triangle {
	u := b.Sub(a)
	v := c.Sub(a)
	return &triangle{
		A:        a,
		U:        u,
		V:        v,
		UnitNorm: u.cross(v).Unit(),
		DotUV:    u.Dot(v),
		DotUU:    u.Dot(u),
		DotVV:    v.Dot(v),
	}
}

func (t *triangle) intersect(r ray) (intersection, bool) {
	// Check if there's a hit with the plane.
	h := t.UnitNorm.Dot(t.A.Sub(r.start)) / t.UnitNorm.Dot(r.dir)
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

	w := r.at(h).Sub(t.A)
	dotWV := w.Dot(t.V)
	dotWU := w.Dot(t.U)
	alpha := t.DotUV*dotWV - t.DotVV*dotWU
	beta := t.DotUV*dotWU - t.DotUU*dotWV
	denom := t.DotUV*t.DotUV - t.DotUU*t.DotVV
	alpha /= denom
	beta /= denom

	if alpha < 0 || beta < 0 || alpha+beta > 1 {
		return intersection{}, false
	}
	return intersection{
		unitNormal: t.UnitNorm,
		distance:   h,
	}, true
}

func (t *triangle) bound() (Vector, Vector) {
	b := t.A.Add(t.U)
	c := t.A.Add(t.V)
	min := t.A.Min(b.Min(c)).addULPs(-ulpFudgeFactor)
	max := t.A.Max(b.Max(c)).addULPs(+ulpFudgeFactor)
	return min, max
}

func (t *triangle) translate(v Vector) {
	a := t.A.Add(v)
	b := t.U.Add(t.A).Add(v)
	c := t.V.Add(t.A).Add(v)
	*t = *newTriangle(a, b, c)
}

func (t *triangle) rotate(v Vector, rads float64) {
	a := t.A.rotate(v, rads)
	b := t.U.Add(t.A).rotate(v, rads)
	c := t.V.Add(t.A).rotate(v, rads)
	*t = *newTriangle(a, b, c)
}

func (t *triangle) scale(f float64) {
	a := t.A.Scale(f)
	b := t.U.Add(t.A).Scale(f)
	c := t.V.Add(t.A).Scale(f)
	*t = *newTriangle(a, b, c)
}

type alignedBox struct {
	Max Vector `json:"max"`
	Min Vector `json:"min"`
}

func (a *alignedBox) String() string {
	return fmt.Sprintf("Type=alignedBox Min=%v Max=%v", a.Max, a.Min)
}

func newAlignedBox(corner1, corner2 Vector) surface {
	return &alignedBox{
		Max: corner1.Min(corner2),
		Min: corner1.Max(corner2),
	}
}

func (b *alignedBox) intersect(r ray) (intersection, bool) {
	tx1 := (b.Max.X - r.start.X) / r.dir.X
	tx2 := (b.Min.X - r.start.X) / r.dir.X
	ty1 := (b.Max.Y - r.start.Y) / r.dir.Y
	ty2 := (b.Min.Y - r.start.Y) / r.dir.Y
	tz1 := (b.Max.Z - r.start.Z) / r.dir.Z
	tz2 := (b.Min.Z - r.start.Z) / r.dir.Z

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
	return b.Max, b.Min
}

func (b *alignedBox) translate(v Vector) {
	b.Max = b.Max.Add(v)
	b.Min = b.Min.Add(v)
}

func (b *alignedBox) rotate(Vector, float64) {
	panic("cannot rotate aligned box")
}

func (b *alignedBox) scale(f float64) {
	b.Max = b.Max.Scale(f)
	b.Min = b.Min.Scale(f)
}

type sphere struct {
	Center Vector  `json:"center"`
	Radius float64 `json:"radius"`
}

func (s *sphere) String() string {
	return fmt.Sprintf("Type=sphere C=%v R=%v", s.Center, s.Radius)
}

func (s *sphere) intersect(r ray) (intersection, bool) {

	// Get coefficients to a.x^2 + b.x + c = 0
	emc := r.start.Sub(s.Center)
	a := r.dir.LengthSq()
	b := 2 * emc.Dot(r.dir)
	c := emc.LengthSq() - s.Radius*s.Radius

	// Find discriminant b*b - 4*a*c
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
		unitNormal: r.at(t).Sub(s.Center).Unit(),
		distance:   t,
	}, t > 0
}

func (s *sphere) bound() (Vector, Vector) {
	r := Vect(s.Radius, s.Radius, s.Radius)
	min, max := s.Center.Sub(r), s.Center.Add(r)
	return min.addULPs(-ulpFudgeFactor), max.addULPs(ulpFudgeFactor)
}

func (s *sphere) translate(v Vector) {
	s.Center = s.Center.Add(v)
}

func (s *sphere) rotate(v Vector, rads float64) {
	// NO-OP
}

func (s *sphere) scale(f float64) {
	s.Radius *= f
	s.Center = s.Center.Scale(f)
}

type alignXSquare struct {
	X  float64 `json:"x"`
	Y1 float64 `json:"y_1"`
	Y2 float64 `json:"y_2"`
	Z1 float64 `json:"z_1"`
	Z2 float64 `json:"z_2"`
}

func (a *alignXSquare) String() string {
	return fmt.Sprintf("Type=alignXSquare X=%v Y1=%v Y2=%v Z1=%v Z2=%v",
		a.X, a.Y1, a.Y2, a.Z1, a.Z2)
}

func (s *alignXSquare) intersect(r ray) (intersection, bool) {
	t := (s.X - r.start.X) / r.dir.X
	hit := r.at(t)
	return intersection{Vect(+1, 0, 0), t},
		t > 0 && hit.Y > s.Y1 && hit.Y < s.Y2 && hit.Z > s.Z1 && hit.Z < s.Z2
}

func (s *alignXSquare) bound() (Vector, Vector) {
	return Vect(s.X, s.Y1, s.Z1), Vect(s.X, s.Y2, s.Z2)
}

func (s *alignXSquare) translate(v Vector) {
	s.X += v.X
	s.Y1 += v.Y
	s.Y2 += v.Y
	s.Z1 += v.Z
	s.Z2 += v.Z
}

func (a *alignXSquare) rotate(Vector, float64) {
	panic("cannot rotate aligned square")
}

func (a *alignXSquare) scale(f float64) {
	a.X *= f
	a.Y1 *= f
	a.Y2 *= f
	a.Z1 *= f
	a.Z2 *= f
}

type alignYSquare struct {
	X1 float64 `json:"x_1"`
	X2 float64 `json:"x_2"`
	Y  float64 `json:"y"`
	Z1 float64 `json:"z_1"`
	Z2 float64 `json:"z_2"`
}

func (a *alignYSquare) String() string {
	return fmt.Sprintf("Type=alignYSquare X1=%v X2=%v Y=%v Z1=%v Z2=%v",
		a.X1, a.X2, a.Y, a.Z1, a.Z2)
}

func (s *alignYSquare) intersect(r ray) (intersection, bool) {
	t := (s.Y - r.start.Y) / r.dir.Y
	hit := r.at(t)
	return intersection{Vect(0, +1, 0), t},
		t > 0 && hit.X > s.X1 && hit.X < s.X2 && hit.Z > s.Z1 && hit.Z < s.Z2
}

func (s *alignYSquare) bound() (Vector, Vector) {
	return Vect(s.X1, s.Y, s.Z1), Vect(s.X2, s.Y, s.Z2)
}

func (s *alignYSquare) translate(v Vector) {
	s.X1 += v.X
	s.X2 += v.X
	s.Y += v.Y
	s.Z1 += v.Z
	s.Z2 += v.Z
}

func (a *alignYSquare) rotate(Vector, float64) {
	panic("cannot rotate aligned square")
}

func (a *alignYSquare) scale(f float64) {
	a.X1 *= f
	a.X2 *= f
	a.Y *= f
	a.Z1 *= f
	a.Z2 *= f
}

type alignZSquare struct {
	X1 float64 `json:"x_1"`
	X2 float64 `json:"x_2"`
	Y1 float64 `json:"y_1"`
	Y2 float64 `json:"y_2"`
	Z  float64 `json:"z"`
}

func (a *alignZSquare) String() string {
	return fmt.Sprintf("Type=alignZSquare X1=%v X2=%v Y1=%v Y2=%v Z=%v",
		a.X1, a.X2, a.Y1, a.Y2, a.Z)
}

func (s *alignZSquare) intersect(r ray) (intersection, bool) {
	t := (s.Z - r.start.Z) / r.dir.Z
	hit := r.at(t)
	return intersection{Vect(0, 0, +1), t},
		t > 0 && hit.X > s.X1 && hit.X < s.X2 && hit.Y > s.Y1 && hit.Y < s.Y2
}

func (s *alignZSquare) bound() (Vector, Vector) {
	return Vect(s.X1, s.Y1, s.Z), Vect(s.X2, s.Y2, s.Z)
}

func (s *alignZSquare) translate(v Vector) {
	s.X1 += v.X
	s.X2 += v.X
	s.Y1 += v.Y
	s.Y2 += v.Y
	s.Z += v.Z
}

func (a *alignZSquare) rotate(Vector, float64) {
	panic("cannot rotate aligned square")
}

func (a *alignZSquare) scale(f float64) {
	a.X1 *= f
	a.X2 *= f
	a.Y1 *= f
	a.Y2 *= f
	a.Z *= f
}

type disc struct {
	Center   Vector  `json:"center"`
	RadiusSq float64 `json:"radius_sq"`
	UnitNorm Vector  `json:"unit_norm"`
}

func (d *disc) intersect(r ray) (intersection, bool) {
	h := d.UnitNorm.Dot(d.Center.Sub(r.start)) / d.UnitNorm.Dot(r.dir)
	if h <= 0 {
		// Hit was behind the camera.
		return intersection{}, false
	}
	hitLoc := r.at(h)
	if hitLoc.Sub(d.Center).LengthSq() > d.RadiusSq {
		return intersection{}, false
	}
	return intersection{
		unitNormal: d.UnitNorm,
		distance:   h,
	}, true
}

func (d *disc) bound() (Vector, Vector) {
	n := d.UnitNorm
	offset := discBoundOffset(n, math.Sqrt(d.RadiusSq))
	return d.Center.Sub(offset), d.Center.Add(offset)
}

func discBoundOffset(n Vector, r float64) Vector {
	assertUnit(n)
	return Vect(n.x0().Length(), n.y0().Length(), n.z0().Length()).Scale(r)
}

func (d *disc) translate(v Vector) {
	d.Center = d.Center.Add(v)
}

func (d *disc) rotate(u Vector, rads float64) {
	d.Center = d.Center.rotate(u, rads)
	d.UnitNorm = d.Center.rotate(u, rads)
}

func (d *disc) scale(s float64) {
	d.Center = d.Center.Scale(s)
	d.RadiusSq *= s * s
}

type pipe struct {
	C1 Vector  `json:"c_1"` // endpoint 1
	C2 Vector  `json:"c_2"` // endpoint 2
	R  float64 `json:"r"`
}

func (p *pipe) String() string {
	return fmt.Sprintf("Type=pipe r=%v c1=%v c2=%v", p.R, p.C1, p.C2)
}

func (p *pipe) intersect(r ray) (intersection, bool) {
	h := p.C2.Sub(p.C1).Unit()
	dCrossH := r.dir.cross(h)
	emc := r.start.Sub(p.C1)
	emcCrossH := emc.cross(h)
	x1, x2 := solveQuadraticEqn(
		dCrossH.LengthSq(),
		2*dCrossH.Dot(emcCrossH),
		emcCrossH.LengthSq()-p.R*p.R,
	)

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	for _, x := range [...]float64{x1, x2} {
		if x <= 0 {
			continue
		}
		hitAt := r.at(x)
		s := hitAt.Sub(p.C1).Dot(h)
		if s < 0 || s*s > p.C2.Sub(p.C1).LengthSq() {
			continue
		}
		return intersection{
			unitNormal: hitAt.Sub(p.C1).rej(h).Unit(),
			distance:   x,
		}, true
	}
	return intersection{}, false
}

func (p *pipe) bound() (Vector, Vector) {
	h := p.C2.Sub(p.C1).Unit()
	offset := discBoundOffset(h, p.R)
	return p.C1.Min(p.C2).Sub(offset), p.C1.Max(p.C2).Add(offset)
}

func (p *pipe) translate(v Vector) {
	p.C1 = p.C1.Add(v)
	p.C2 = p.C2.Add(v)
}

func (p *pipe) rotate(v Vector, rads float64) {
	p.C1 = p.C1.rotate(v, rads)
	p.C2 = p.C2.rotate(v, rads)
}

func (p *pipe) scale(s float64) {
	p.C1 = p.C1.Scale(s)
	p.C2 = p.C2.Scale(s)
}

func solveQuadraticEqn(a, b, c float64) (float64, float64) {
	disc := b*b - 4*a*c
	if disc < 0 {
		return -1, -1
	}

	// Find x1 and x2 using a numerically stable algorithm.
	var signOfB float64
	signOfB = math.Copysign(1.0, b)
	q := -0.5 * (b + signOfB*math.Sqrt(disc))
	return q / a, c / q
}
