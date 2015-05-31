package grayt

import (
	"math"
	"math/rand"
)

func NewSphere(centre Vect, radius float64) Surface {
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

	return Intersection{r.At(t).Sub(s.centre).Unit(), t}, t > 0
}

func (s *sphere) BoundingBox() (min, max Vect) {
	r := Vect{s.radius, s.radius, s.radius}
	min = s.centre.Sub(r)
	max = s.centre.Add(r)
	return
}

func (s *sphere) Sample() Vect {
	rnd := Vect{rand.NormFloat64(), rand.NormFloat64(), rand.NormFloat64()}
	return rnd.Extended(s.radius / rnd.Length()).Add(s.centre)
}
