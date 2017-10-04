package grayt

import "math"

func unitSphere(e, d Vector) (hloc, norm Vector) {
	a := d.LengthSq()
	b := 2 * e.Dot(d)
	c := e.LengthSq() - 1
	t := solveQuad(a, b, c)
	if t <= 0 {
		return
	}
	hloc = e.Add(d.Scale(t))
	norm = hloc
	return
}

// solveQuad finds the smallest positive solution (if there is one) to a*x*x +
// b*x + c == 0, otherwise returns a negative number.
func solveQuad(a, b, c float64) float64 {
	disc := b*b - 4*a*c
	if disc < 0 {
		return -1
	}

	// Find x1 and x2 using a numerically stable algorithm.
	var signOfB float64
	signOfB = math.Copysign(1.0, b)
	q := -0.5 * (b + signOfB*math.Sqrt(disc))
	x1 := q / a
	x2 := c / q

	if x1 > 0 && x2 > 0 {
		// Both are positive, so take the smaller one.
		return math.Min(x1, x2)
	} else {
		// At least one is negative, take the larger one (which is either
		// negative or positive).
		return math.Max(x1, x2)
	}
}
