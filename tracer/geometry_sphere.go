package tracer

import (
	"math"

	"github.com/peterstace/grayt/ray"
	"github.com/peterstace/grayt/vect"
)

type sphere struct {
	centre vect.V
	radius float64
}

// NewSphere creates a sphere with the given centre and radius.
func NewSphere(centre vect.V, radius float64) Geometry {
	return &sphere{centre: centre, radius: radius}
}

func (s *sphere) intersect(r ray.Ray) (hitRec, bool) {
	t := s.t(r)
	if t <= 0.0 {
		return hitRec{}, false
	}
	return hitRec{t: t, n: vect.Sub(r.At(t), s.centre)}, true
}

func (s *sphere) t(r ray.Ray) float64 {

	// Get coeficients to a.x^2 + b.x + c = 0
	emc := vect.Sub(r.Start, s.centre)
	a := r.Dir.Length2()
	b := 2 * vect.Dot(emc, r.Dir)
	c := emc.Length2() - s.radius*s.radius

	// Find discrimenant b*b - 4*a*c
	disc := b*b - 4*a*c
	if disc < 0 {
		return -1.0
	}

	// solve for x1 and x2
	var signOfB float64
	if b > 0 {
		signOfB = 1.0
	} else if b < 0 {
		signOfB = -1.0
	}
	q := -0.5 * (b * signOfB * math.Sqrt(disc))
	x1 := q / a
	x2 := c / q

	if x1 > 0 && x2 > 0 {
		// Both are positive, so take the smaller one.
		return math.Min(x1, x2)
	}
	// At least one is negative, take the larger one (which is either
	// negative or positive).
	return math.Max(x1, x2)
}
