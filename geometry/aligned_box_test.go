package grayt

import "testing"

func TestAlignedBoxIntersect(t *testing.T) {
	b := AlignedBox{
		Corner1: Vector{-1, -2, -3},
		Corner2: Vector{1, 2, 3},
	}.MakeSurfaces()[0]
	tests := []struct {
		r Ray
		t float64
		n Vector
	}{
		{
			r: Ray{Start: Vector{-5, 0, 0}, Dir: Vector{1, 0, 0}},
			t: 4.0,
			n: Vector{-1, 0, 0},
		},
		{
			r: Ray{Start: Vector{0, -5, 0}, Dir: Vector{0, 1, 0}},
			t: 3.0,
			n: Vector{0, -1, 0},
		},
		{
			r: Ray{Start: Vector{0, 0, -5}, Dir: Vector{0, 0, 1}},
			t: 2.0,
			n: Vector{0, 0, -1},
		},
		{
			r: Ray{Start: Vector{5, 0, 0}, Dir: Vector{-1, 0, 0}},
			t: 4.0,
			n: Vector{1, 0, 0},
		},
		{
			r: Ray{Start: Vector{0, 5, 0}, Dir: Vector{0, -1, 0}},
			t: 3.0,
			n: Vector{0, 1, 0},
		},
		{
			r: Ray{Start: Vector{0, 0, 5}, Dir: Vector{0, 0, -1}},
			t: 2.0,
			n: Vector{0, 0, 1},
		},

		{
			r: Ray{Start: Vector{0, 0, 0}, Dir: Vector{1, 0, 0}},
			t: 1.0,
			n: Vector{1, 0, 0},
		},
		{
			r: Ray{Start: Vector{0, 0, 0}, Dir: Vector{0, 1, 0}},
			t: 2.0,
			n: Vector{0, 1, 0},
		},
		{
			r: Ray{Start: Vector{0, 0, 0}, Dir: Vector{0, 0, 1}},
			t: 3.0,
			n: Vector{0, 0, 1},
		},
		{
			r: Ray{Start: Vector{0, 0, 0}, Dir: Vector{-1, 0, 0}},
			t: 1.0,
			n: Vector{-1, 0, 0},
		},
		{
			r: Ray{Start: Vector{0, 0, 0}, Dir: Vector{0, -1, 0}},
			t: 2.0,
			n: Vector{0, -1, 0},
		},
		{
			r: Ray{Start: Vector{0, 0, 0}, Dir: Vector{0, 0, -1}},
			t: 3.0,
			n: Vector{0, 0, -1},
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
