package grayt

import "math"

type Sphere struct {
	Center Vect
	Radius float64
}

func (s *Sphere) HitSurface(r ray) (float64, Vect) {
	if t := s.calcT(r); t >= 0 {
		return t, Sub(r.at(t), s.Center)
	}
	return -1.0, Vect{}
}

func (s *Sphere) calcT(r ray) float64 {

	// get coeficients to a.x^2 + b.x + c = 0
	emc := Sub(r.start, s.Center)
	a := r.dir.Length2()
	b := 2 * Dot(emc, r.dir)
	c := emc.Length2() - s.Radius*s.Radius

	// find discrimenant b*b - 4*a*c
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
	} else {
		signOfB = 0.0
	}
	q := -0.5 * (b + signOfB*math.Sqrt(disc))
	x1 := q / a
	x2 := c / q

	// Get the smallest positive solution (or return negative if there is no positive solution).
	if x1 > 0 && x2 > 0 {
		// Both are positive, we want the smaller one.
		if x1 < x2 {
			return x1
		} else {
			return x2
		}
	} else {
		// At least one is negative, take the larger one (which is either positive or negative).
		if x1 > x2 {
			return x1
		} else {
			return x2
		}
	}
}
