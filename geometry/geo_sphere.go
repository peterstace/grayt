package geometry

import (
	"math"

	"github.com/peterstace/grayt"
)

type Sphere struct {
	C grayt.Vector
	R float64
}

func (s Sphere) MakeSurfaces() []grayt.Surface {
	return []grayt.Surface{&sphere{centre: s.C, radius: s.R}}
}

type sphere struct {
	centre grayt.Vector
	radius float64
}

func (s *sphere) Intersect(r grayt.Ray) (grayt.Intersection, bool) {

	// Get coeficients to a.x^2 + b.x + c = 0
	emc := r.Start.Sub(s.centre)
	a := r.Dir.LengthSq()
	b := 2 * emc.Dot(r.Dir)
	c := emc.LengthSq() - s.radius*s.radius

	// Find discrimenant b*b - 4*a*c
	disc := b*b - 4*a*c
	if disc < 0 {
		return grayt.Intersection{}, false
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

	return grayt.Intersection{r.At(t).Sub(s.centre).Unit(), t}, t > 0
}
