package grayt

import "math"

type intersection struct {
	distance   float64
	unitNormal Vect
}

type geometry interface {
	intersect(Ray) (intersection, bool)
}

type plane struct {
	unitNormal Vect // Unit normal out of the plane.
	anchor     Vect // Any point on the plane.
}

func (p *plane) intersect(r Ray) (intersection, bool) {
	t := p.unitNormal.Dot(p.anchor.Sub(r.Start)) / p.unitNormal.Dot(r.Dir)
	return intersection{distance: t, unitNormal: p.unitNormal}, t > 0
}

type sphere struct {
	centre Vect
	radius float64
}

func (s *sphere) intersect(r Ray) (intersection, bool) {

	// Get coeficients to a.x^2 + b.x + c = 0
	emc := r.Start.Sub(s.centre)
	a := r.Dir.Length2()
	b := 2 * emc.Dot(r.Dir)
	c := emc.Length2() - s.radius*s.radius

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

	return intersection{distance: t, unitNormal: r.At(t).Sub(s.centre).Unit()}, t > 0
}
