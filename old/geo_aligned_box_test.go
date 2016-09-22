package grayt

import "testing"

func TestAlignedBoxIntersect(t *testing.T) {
	b := AlignedBox(
		Vect(-1, -2, -3),
		Vect(1, 2, 3),
	)
	tests := []struct {
		r Ray
		t float64
		n Vector
	}{
		{
			r: Ray{Start: Vect(-5, 0, 0), Dir: Vect(1, 0, 0)},
			t: 4.0,
			n: Vect(-1, 0, 0),
		},
		{
			r: Ray{Start: Vect(0, -5, 0), Dir: Vect(0, 1, 0)},
			t: 3.0,
			n: Vect(0, -1, 0),
		},
		{
			r: Ray{Start: Vect(0, 0, -5), Dir: Vect(0, 0, 1)},
			t: 2.0,
			n: Vect(0, 0, -1),
		},
		{
			r: Ray{Start: Vect(5, 0, 0), Dir: Vect(-1, 0, 0)},
			t: 4.0,
			n: Vect(1, 0, 0),
		},
		{
			r: Ray{Start: Vect(0, 5, 0), Dir: Vect(0, -1, 0)},
			t: 3.0,
			n: Vect(0, 1, 0),
		},
		{
			r: Ray{Start: Vect(0, 0, 5), Dir: Vect(0, 0, -1)},
			t: 2.0,
			n: Vect(0, 0, 1),
		},

		{
			r: Ray{Start: Vect(0, 0, 0), Dir: Vect(1, 0, 0)},
			t: 1.0,
			n: Vect(1, 0, 0),
		},
		{
			r: Ray{Start: Vect(0, 0, 0), Dir: Vect(0, 1, 0)},
			t: 2.0,
			n: Vect(0, 1, 0),
		},
		{
			r: Ray{Start: Vect(0, 0, 0), Dir: Vect(0, 0, 1)},
			t: 3.0,
			n: Vect(0, 0, 1),
		},
		{
			r: Ray{Start: Vect(0, 0, 0), Dir: Vect(-1, 0, 0)},
			t: 1.0,
			n: Vect(-1, 0, 0),
		},
		{
			r: Ray{Start: Vect(0, 0, 0), Dir: Vect(0, -1, 0)},
			t: 2.0,
			n: Vect(0, -1, 0),
		},
		{
			r: Ray{Start: Vect(0, 0, 0), Dir: Vect(0, 0, -1)},
			t: 3.0,
			n: Vect(0, 0, -1),
		},
	}
	for i, test := range tests {
		isect, hit := b.Intersect(test.r)
		if !hit {
			t.Errorf("%d: didn't hit, %v", i, test)
		}
		if isect.UnitNormal != test.n {
			t.Errorf("%d: want %v, got %v", i, test.n, isect.UnitNormal)
		}
		if isect.Distance != test.t {
			t.Errorf("%d: want %g, got %g", i, test.t, isect.Distance)
		}
	}
}
