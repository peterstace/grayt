package geometry

import (
	"testing"

	"github.com/peterstace/grayt"
)

func TestAlignedBoxIntersect(t *testing.T) {
	b := AlignedBox{
		Corner1: grayt.Vector{-1, -2, -3},
		Corner2: grayt.Vector{1, 2, 3},
	}.MakeSurfaces()[0]
	tests := []struct {
		r grayt.Ray
		t float64
		n grayt.Vector
	}{
		{
			r: grayt.Ray{Start: grayt.Vector{-5, 0, 0}, Dir: grayt.Vector{1, 0, 0}},
			t: 4.0,
			n: grayt.Vector{-1, 0, 0},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{0, -5, 0}, Dir: grayt.Vector{0, 1, 0}},
			t: 3.0,
			n: grayt.Vector{0, -1, 0},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{0, 0, -5}, Dir: grayt.Vector{0, 0, 1}},
			t: 2.0,
			n: grayt.Vector{0, 0, -1},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{5, 0, 0}, Dir: grayt.Vector{-1, 0, 0}},
			t: 4.0,
			n: grayt.Vector{1, 0, 0},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{0, 5, 0}, Dir: grayt.Vector{0, -1, 0}},
			t: 3.0,
			n: grayt.Vector{0, 1, 0},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{0, 0, 5}, Dir: grayt.Vector{0, 0, -1}},
			t: 2.0,
			n: grayt.Vector{0, 0, 1},
		},

		{
			r: grayt.Ray{Start: grayt.Vector{0, 0, 0}, Dir: grayt.Vector{1, 0, 0}},
			t: 1.0,
			n: grayt.Vector{1, 0, 0},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{0, 0, 0}, Dir: grayt.Vector{0, 1, 0}},
			t: 2.0,
			n: grayt.Vector{0, 1, 0},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{0, 0, 0}, Dir: grayt.Vector{0, 0, 1}},
			t: 3.0,
			n: grayt.Vector{0, 0, 1},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{0, 0, 0}, Dir: grayt.Vector{-1, 0, 0}},
			t: 1.0,
			n: grayt.Vector{-1, 0, 0},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{0, 0, 0}, Dir: grayt.Vector{0, -1, 0}},
			t: 2.0,
			n: grayt.Vector{0, -1, 0},
		},
		{
			r: grayt.Ray{Start: grayt.Vector{0, 0, 0}, Dir: grayt.Vector{0, 0, -1}},
			t: 3.0,
			n: grayt.Vector{0, 0, -1},
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
