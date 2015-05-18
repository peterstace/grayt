package grayt

import "math"

// Intersection between some geometry and a ray.
type Intersection struct {
	Distance   float64 // Distance along the ray where the intersection occurred.
	UnitNormal Vect    // Unit normal (pointing 'away' from the geometry, not 'into' it).
}

// Geometry implementations represent surfaces that can be intersected with.
type Geometry interface {

	// Intersect finds the intersection (if it exists) between a ray and the
	// geometry.
	Intersect(Ray) (Intersection, bool)
}

func NewPlane(normal, anchor Vect) Geometry {
	return &plane{
		unitNormal: normal.Unit(),
		anchor:     anchor,
	}
}

type plane struct {
	unitNormal Vect // Unit normal out of the plane.
	anchor     Vect // Any point on the plane.
}

func (p *plane) Intersect(r Ray) (Intersection, bool) {
	t := p.unitNormal.Dot(p.anchor.Sub(r.Start)) / p.unitNormal.Dot(r.Dir)
	return Intersection{Distance: t, UnitNormal: p.unitNormal}, t > 0
}

func NewSphere(centre Vect, radius float64) Geometry {
	return &sphere{
		centre: centre,
		radius: radius,
	}
}

type sphere struct {
	centre Vect
	radius float64
}

func (s *sphere) Intersect(r Ray) (Intersection, bool) {

	// Get coeficients to a.x^2 + b.x + c = 0
	emc := r.Start.Sub(s.centre)
	a := r.Dir.Length2()
	b := 2 * emc.Dot(r.Dir)
	c := emc.Length2() - s.radius*s.radius

	// Find discrimenant b*b - 4*a*c
	disc := b*b - 4*a*c
	if disc < 0 {
		return Intersection{}, false
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

	return Intersection{Distance: t, UnitNormal: r.At(t).Sub(s.centre).Unit()}, t > 0
}
