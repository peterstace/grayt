package trace

import (
	"testing"

	"github.com/peterstace/grayt/xmath"
)

func TestAlignedBoxIntersection(t *testing.T) {
	box := alignedBox{
		xmath.Vect(0.9, 1.0, -0.9),
		xmath.Vect(0.1, 0.999, -0.1),
	}
	r := xmath.Ray{
		xmath.Vect(0.18325558497391392, 0.999012404240046, -0.9999999999999999),
		xmath.Vect(-0.5720297174078126, 0.007828964018125696, 0.8201955313976911),
	}
	intersection, hit := box.intersect(r)
	if !hit {
		t.Error("should have hit")
	}
	if intersection.distance != 0.12192214681978378 {
		t.Errorf("wrong distance")
	}
	if intersection.unitNormal != xmath.Vect(0, 0, -1) {
		t.Errorf("wrong normal")
	}
}

func BenchmarkTriangleIntersect(b *testing.B) {

	t := newTriangle(
		xmath.Vect(+0, 0, +0),
		xmath.Vect(-1, 1, -1),
		xmath.Vect(+1, 1, -1),
	)
	r := xmath.Ray{Start: xmath.Vect(0, 0.5, 5), Dir: xmath.Vect(0, 0, -1)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.intersect(r)
	}
}
