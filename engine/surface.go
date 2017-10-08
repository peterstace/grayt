package engine

import (
	"fmt"
	"math"
)

type surface struct {
	transform matrix4
	primitive func(e, d vect3) (n, h vect3, ok bool)
}

func (s *surface) intersect(e, d vect3) (n, h vect3, ok bool) {

	transInv, ok := s.transform.inv()
	if !ok {
		panic("could not invert transformation matrix")
	}

	fmt.Println("  trans", s.transform)
	fmt.Println("  transInv", transInv)
	ePrime := transInv.mulv(e.extend(1)).truncPoint()
	dPrime := transInv.mulv(d.extend(0)).truncVect().unit()

	fmt.Println("  ePrime", ePrime)
	fmt.Println("  dPrime", dPrime)

	var nPrime, hPrime vect3
	nPrime, hPrime, ok = triangle(ePrime, dPrime)
	if !ok {
		return
	}

	n = transInv.transpose().mulv(nPrime.extend(0)).truncVect()
	h = s.transform.mulv(hPrime.extend(1)).truncPoint()
	return
}

var invSqrt3 = 1 / math.Sqrt(3)

func triangle(e, d vect3) (n, h vect3, ok bool) {
	assertUnit(d)
	const third = 1.0 / 3.0
	n = vect3{invSqrt3, invSqrt3, invSqrt3}
	t := vect3{third, third, third}.sub(e).dot(n) / d.dot(n)
	fmt.Println("t", t)
	if t < 0 {
		return
	}
	h = e.add(d.scale(t))
	fmt.Println("h", h)
	if h[0] < 0 || h[1] < 0 || h[2] < 0 {
		return
	}
	ok = true
	return
}
