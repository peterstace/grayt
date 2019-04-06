package trace

import (
	"fmt"
	"math"

	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/xmath"
)

const ulpFudgeFactor = 50

type surface interface {
	intersect(r xmath.Ray) (intersection, bool)
	bound() (xmath.Vector, xmath.Vector)
	translate(xmath.Vector)
	rotate(xmath.Vector, float64)
	scale(float64)
}

type material struct {
	Colour    colour.Colour `json:"colour"`
	Emittance float64       `json:"emittance"`
	Mirror    bool          `json:"mirror"`
}

type object struct {
	Surface  surface  `json:"surface"`
	Material material `json:"material"`
}

func (o object) String() string {
	return fmt.Sprintf("Surface={%v} Material={%v}", o.Surface, o.Material)
}

type intersection struct {
	unitNormal xmath.Vector
	distance   float64
}

type triangle struct {
	A        xmath.Vector `json:"a"` // Corner A
	U        xmath.Vector `json:"u"` // A to B
	V        xmath.Vector `json:"v"` // A to C
	UnitNorm xmath.Vector `json:"unit_norm"`
	// Precomputed dot products:
	DotUV float64 `json:"dot_uv"`
	DotUU float64 `json:"dot_uu"`
	DotVV float64 `json:"dot_vv"`
}

func (t *triangle) String() string {
	return fmt.Sprintf("Type=triangle A=%v B=%v C=%v", t.A, t.A.Add(t.U), t.A.Add(t.V))
}

func newTriangle(a, b, c xmath.Vector) *triangle {
	u := b.Sub(a)
	v := c.Sub(a)
	return &triangle{
		A:        a,
		U:        u,
		V:        v,
		UnitNorm: u.Cross(v).Unit(),
		DotUV:    u.Dot(v),
		DotUU:    u.Dot(u),
		DotVV:    v.Dot(v),
	}
}

func (t *triangle) intersect(r xmath.Ray) (intersection, bool) {
	// Check if there's a hit with the plane.
	h := t.UnitNorm.Dot(t.A.Sub(r.Start)) / t.UnitNorm.Dot(r.Dir)
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

	w := r.At(h).Sub(t.A)
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

func (t *triangle) bound() (xmath.Vector, xmath.Vector) {
	b := t.A.Add(t.U)
	c := t.A.Add(t.V)
	min := t.A.Min(b.Min(c)).AddULPs(-ulpFudgeFactor)
	max := t.A.Max(b.Max(c)).AddULPs(+ulpFudgeFactor)
	return min, max
}

func (t *triangle) translate(v xmath.Vector) {
	a := t.A.Add(v)
	b := t.U.Add(t.A).Add(v)
	c := t.V.Add(t.A).Add(v)
	*t = *newTriangle(a, b, c)
}

func (t *triangle) rotate(v xmath.Vector, rads float64) {
	a := t.A.Rotate(v, rads)
	b := t.U.Add(t.A).Rotate(v, rads)
	c := t.V.Add(t.A).Rotate(v, rads)
	*t = *newTriangle(a, b, c)
}

func (t *triangle) scale(f float64) {
	a := t.A.Scale(f)
	b := t.U.Add(t.A).Scale(f)
	c := t.V.Add(t.A).Scale(f)
	*t = *newTriangle(a, b, c)
}

type alignedBox struct {
	Max xmath.Vector `json:"max"`
	Min xmath.Vector `json:"min"`
}

func (a *alignedBox) String() string {
	return fmt.Sprintf("Type=alignedBox Min=%v Max=%v", a.Max, a.Min)
}

func newAlignedBox(corner1, corner2 xmath.Vector) surface {
	return &alignedBox{
		Max: corner1.Min(corner2),
		Min: corner1.Max(corner2),
	}
}

func (b *alignedBox) intersect(r xmath.Ray) (intersection, bool) {
	tx1 := (b.Max.X - r.Start.X) / r.Dir.X
	tx2 := (b.Min.X - r.Start.X) / r.Dir.X
	ty1 := (b.Max.Y - r.Start.Y) / r.Dir.Y
	ty2 := (b.Min.Y - r.Start.Y) / r.Dir.Y
	tz1 := (b.Max.Z - r.Start.Z) / r.Dir.Z
	tz2 := (b.Min.Z - r.Start.Z) / r.Dir.Z

	tmin, tmax := math.Inf(-1), math.Inf(+1)
	var nMin xmath.Vector
	var nMax xmath.Vector

	if math.Min(tx1, tx2) > tmin {
		if tx1 < tx2 {
			tmin = tx1
			nMin = xmath.Vect(-1, 0, 0)
		} else {
			tmin = tx2
			nMin = xmath.Vect(1, 0, 0)
		}
	}
	if math.Max(tx1, tx2) < tmax {
		if tx1 > tx2 {
			tmax = tx1
			nMax = xmath.Vect(-1, 0, 0)
		} else {
			tmax = tx2
			nMax = xmath.Vect(1, 0, 0)
		}
	}

	if math.Min(ty1, ty2) > tmin {
		if ty1 < ty2 {
			tmin = ty1
			nMin = xmath.Vect(0, -1, 0)
		} else {
			tmin = ty2
			nMin = xmath.Vect(0, 1, 0)
		}
	}
	if math.Max(ty1, ty2) < tmax {
		if ty1 > ty2 {
			tmax = ty1
			nMax = xmath.Vect(0, -1, 0)
		} else {
			tmax = ty2
			nMax = xmath.Vect(0, 1, 0)
		}
	}

	if math.Min(tz1, tz2) > tmin {
		if tz1 < tz2 {
			tmin = tz1
			nMin = xmath.Vect(0, 0, -1)
		} else {
			tmin = tz2
			nMin = xmath.Vect(0, 0, 1)
		}
	}
	if math.Max(tz1, tz2) < tmax {
		if tz1 > tz2 {
			tmax = tz1
			nMax = xmath.Vect(0, 0, -1)
		} else {
			tmax = tz2
			nMax = xmath.Vect(0, 0, 1)
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

func (b *alignedBox) bound() (xmath.Vector, xmath.Vector) {
	return b.Max, b.Min
}

func (b *alignedBox) translate(v xmath.Vector) {
	b.Max = b.Max.Add(v)
	b.Min = b.Min.Add(v)
}

func (b *alignedBox) rotate(xmath.Vector, float64) {
	panic("cannot rotate aligned box")
}

func (b *alignedBox) scale(f float64) {
	b.Max = b.Max.Scale(f)
	b.Min = b.Min.Scale(f)
}

type sphere struct {
	Center xmath.Vector `json:"center"`
	Radius float64      `json:"radius"`
}

func (s *sphere) String() string {
	return fmt.Sprintf("Type=sphere C=%v R=%v", s.Center, s.Radius)
}

func (s *sphere) intersect(r xmath.Ray) (intersection, bool) {

	// Get coefficients to a.x^2 + b.x + c = 0
	emc := r.Start.Sub(s.Center)
	a := r.Dir.LengthSq()
	b := 2 * emc.Dot(r.Dir)
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
		unitNormal: r.At(t).Sub(s.Center).Unit(),
		distance:   t,
	}, t > 0
}

func (s *sphere) bound() (xmath.Vector, xmath.Vector) {
	r := xmath.Vect(s.Radius, s.Radius, s.Radius)
	min, max := s.Center.Sub(r), s.Center.Add(r)
	return min.AddULPs(-ulpFudgeFactor), max.AddULPs(ulpFudgeFactor)
}

func (s *sphere) translate(v xmath.Vector) {
	s.Center = s.Center.Add(v)
}

func (s *sphere) rotate(v xmath.Vector, rads float64) {
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

func (s *alignXSquare) intersect(r xmath.Ray) (intersection, bool) {
	t := (s.X - r.Start.X) / r.Dir.X
	hit := r.At(t)
	return intersection{xmath.Vect(+1, 0, 0), t},
		t > 0 && hit.Y > s.Y1 && hit.Y < s.Y2 && hit.Z > s.Z1 && hit.Z < s.Z2
}

func (s *alignXSquare) bound() (xmath.Vector, xmath.Vector) {
	return xmath.Vect(s.X, s.Y1, s.Z1), xmath.Vect(s.X, s.Y2, s.Z2)
}

func (s *alignXSquare) translate(v xmath.Vector) {
	s.X += v.X
	s.Y1 += v.Y
	s.Y2 += v.Y
	s.Z1 += v.Z
	s.Z2 += v.Z
}

func (a *alignXSquare) rotate(xmath.Vector, float64) {
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

func (s *alignYSquare) intersect(r xmath.Ray) (intersection, bool) {
	t := (s.Y - r.Start.Y) / r.Dir.Y
	hit := r.At(t)
	return intersection{xmath.Vect(0, +1, 0), t},
		t > 0 && hit.X > s.X1 && hit.X < s.X2 && hit.Z > s.Z1 && hit.Z < s.Z2
}

func (s *alignYSquare) bound() (xmath.Vector, xmath.Vector) {
	return xmath.Vect(s.X1, s.Y, s.Z1), xmath.Vect(s.X2, s.Y, s.Z2)
}

func (s *alignYSquare) translate(v xmath.Vector) {
	s.X1 += v.X
	s.X2 += v.X
	s.Y += v.Y
	s.Z1 += v.Z
	s.Z2 += v.Z
}

func (a *alignYSquare) rotate(xmath.Vector, float64) {
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

func (s *alignZSquare) intersect(r xmath.Ray) (intersection, bool) {
	t := (s.Z - r.Start.Z) / r.Dir.Z
	hit := r.At(t)
	return intersection{xmath.Vect(0, 0, +1), t},
		t > 0 && hit.X > s.X1 && hit.X < s.X2 && hit.Y > s.Y1 && hit.Y < s.Y2
}

func (s *alignZSquare) bound() (xmath.Vector, xmath.Vector) {
	return xmath.Vect(s.X1, s.Y1, s.Z), xmath.Vect(s.X2, s.Y2, s.Z)
}

func (s *alignZSquare) translate(v xmath.Vector) {
	s.X1 += v.X
	s.X2 += v.X
	s.Y1 += v.Y
	s.Y2 += v.Y
	s.Z += v.Z
}

func (a *alignZSquare) rotate(xmath.Vector, float64) {
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
	Center   xmath.Vector `json:"center"`
	RadiusSq float64      `json:"radius_sq"`
	UnitNorm xmath.Vector `json:"unit_norm"`
}

func (d *disc) intersect(r xmath.Ray) (intersection, bool) {
	h := d.UnitNorm.Dot(d.Center.Sub(r.Start)) / d.UnitNorm.Dot(r.Dir)
	if h <= 0 {
		// Hit was behind the camera.
		return intersection{}, false
	}
	hitLoc := r.At(h)
	if hitLoc.Sub(d.Center).LengthSq() > d.RadiusSq {
		return intersection{}, false
	}
	return intersection{
		unitNormal: d.UnitNorm,
		distance:   h,
	}, true
}

func (d *disc) bound() (xmath.Vector, xmath.Vector) {
	n := d.UnitNorm
	offset := discBoundOffset(n, math.Sqrt(d.RadiusSq))
	return d.Center.Sub(offset), d.Center.Add(offset)
}

func discBoundOffset(n xmath.Vector, r float64) xmath.Vector {
	assertUnit(n)
	return xmath.Vect(n.X0().Length(), n.Y0().Length(), n.Z0().Length()).Scale(r)
}

func (d *disc) translate(v xmath.Vector) {
	d.Center = d.Center.Add(v)
}

func (d *disc) rotate(u xmath.Vector, rads float64) {
	d.Center = d.Center.Rotate(u, rads)
	d.UnitNorm = d.Center.Rotate(u, rads)
}

func (d *disc) scale(s float64) {
	d.Center = d.Center.Scale(s)
	d.RadiusSq *= s * s
}

type pipe struct {
	C1 xmath.Vector `json:"c_1"` // endpoint 1
	C2 xmath.Vector `json:"c_2"` // endpoint 2
	R  float64      `json:"r"`
}

func (p *pipe) String() string {
	return fmt.Sprintf("Type=pipe r=%v c1=%v c2=%v", p.R, p.C1, p.C2)
}

func (p *pipe) intersect(r xmath.Ray) (intersection, bool) {
	h := p.C2.Sub(p.C1).Unit()
	dCrossH := r.Dir.Cross(h)
	emc := r.Start.Sub(p.C1)
	emcCrossH := emc.Cross(h)
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
		hitAt := r.At(x)
		s := hitAt.Sub(p.C1).Dot(h)
		if s < 0 || s*s > p.C2.Sub(p.C1).LengthSq() {
			continue
		}
		return intersection{
			unitNormal: hitAt.Sub(p.C1).Rej(h).Unit(),
			distance:   x,
		}, true
	}
	return intersection{}, false
}

func (p *pipe) bound() (xmath.Vector, xmath.Vector) {
	h := p.C2.Sub(p.C1).Unit()
	offset := discBoundOffset(h, p.R)
	return p.C1.Min(p.C2).Sub(offset), p.C1.Max(p.C2).Add(offset)
}

func (p *pipe) translate(v xmath.Vector) {
	p.C1 = p.C1.Add(v)
	p.C2 = p.C2.Add(v)
}

func (p *pipe) rotate(v xmath.Vector, rads float64) {
	p.C1 = p.C1.Rotate(v, rads)
	p.C2 = p.C2.Rotate(v, rads)
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
