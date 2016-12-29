package grayt

import "testing"

func TestAlignedBoxIntersection(t *testing.T) {
	_ = "breakpoint"

	box := AlignedBox(
		Vect(0.9, 1.0, -0.9),
		Vect(0.1, 0.999, -0.1),
	)[0]
	r := ray{
		Vector{0.18325558497391392, 0.999012404240046, -0.9999999999999999},
		Vector{-0.5720297174078126, 0.007828964018125696, 0.8201955313976911},
	}
	intersection, hit := box.intersect(r)
	if !hit {
		t.Error("should have hit")
	}
	if intersection.distance != 0.12192214681978378 {
		t.Errorf("wrong distance")
	}
	if intersection.unitNormal != Vect(0, 0, -1) {
		t.Errorf("wrong normal")
	}
}
