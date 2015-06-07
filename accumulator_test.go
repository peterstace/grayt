package grayt

import "testing"

func TestMean(t *testing.T) {
	acc := Accumulator{acc: []Colour{
		Colour{2, 4, 4},
		Colour{4, 4, 10},
		Colour{7, 9, 10},
	}}
	if mean := acc.mean(); mean != 6.0 {
		t.Errorf("got=%f want=6.0", mean)
	}
}
