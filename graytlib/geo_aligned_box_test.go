package graytlib

import "testing"

func TestAlignedBoxIntersect(t *testing.T) {
	b := AlignedBox{
		Corner1: Vect{-1, -2, -3},
		Corner2: Vect{1, 2, 3},
	}.MakeSurfaces()[0]
	tests := []struct {
		r Ray
		t float64
		n Vect
	}{
		{
			r: Ray{Start: Vect{-5, 0, 0}, Dir: Vect{1, 0, 0}},
			t: 4.0,
			n: Vect{-1, 0, 0},
		},
		{
			r: Ray{Start: Vect{0, -5, 0}, Dir: Vect{0, 1, 0}},
			t: 3.0,
			n: Vect{0, -1, 0},
		},
		{
			r: Ray{Start: Vect{0, 0, -5}, Dir: Vect{0, 0, 1}},
			t: 2.0,
			n: Vect{0, 0, -1},
		},
		{
			r: Ray{Start: Vect{5, 0, 0}, Dir: Vect{-1, 0, 0}},
			t: 4.0,
			n: Vect{1, 0, 0},
		},
		{
			r: Ray{Start: Vect{0, 5, 0}, Dir: Vect{0, -1, 0}},
			t: 3.0,
			n: Vect{0, 1, 0},
		},
		{
			r: Ray{Start: Vect{0, 0, 5}, Dir: Vect{0, 0, -1}},
			t: 2.0,
			n: Vect{0, 0, 1},
		},

		{
			r: Ray{Start: Vect{0, 0, 0}, Dir: Vect{1, 0, 0}},
			t: 1.0,
			n: Vect{1, 0, 0},
		},
		{
			r: Ray{Start: Vect{0, 0, 0}, Dir: Vect{0, 1, 0}},
			t: 2.0,
			n: Vect{0, 1, 0},
		},
		{
			r: Ray{Start: Vect{0, 0, 0}, Dir: Vect{0, 0, 1}},
			t: 3.0,
			n: Vect{0, 0, 1},
		},
		{
			r: Ray{Start: Vect{0, 0, 0}, Dir: Vect{-1, 0, 0}},
			t: 1.0,
			n: Vect{-1, 0, 0},
		},
		{
			r: Ray{Start: Vect{0, 0, 0}, Dir: Vect{0, -1, 0}},
			t: 2.0,
			n: Vect{0, -1, 0},
		},
		{
			r: Ray{Start: Vect{0, 0, 0}, Dir: Vect{0, 0, -1}},
			t: 3.0,
			n: Vect{0, 0, -1},
		},
	}
	for _, test := range tests {
		isect, hit := b.Intersect(test.r)
		if !hit {
			t.Errorf("didn't hit, %v", test)
		}
		if isect.UnitNormal != test.n {
			t.Errorf("want %v, got %v", test.n, isect.UnitNormal)
		}
		if isect.Distance != test.t {
			t.Errorf("want %g, got %g", test.t, isect.Distance)
		}
	}
}
