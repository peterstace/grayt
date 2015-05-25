package grayt

import "testing"

func TestDistribution(t *testing.T) {
	acc := Accumulator{acc: []float64{2, 4, 4, 4, 5, 5, 7, 9}}
	if mean := acc.mean(); mean != 5.0 {
		t.Errorf("got=%f want=5.0", mean)
	}
}
