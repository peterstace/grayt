package grayt

import "testing"

func TestDistribution(t *testing.T) {
	acc := Accumulator{acc: []float64{2, 4, 4, 4, 5, 5, 7, 9}}
	mean, stddev := acc.distribution()
	if mean != 5.0 {
		t.Errorf("got=%f want=5.0", mean)
	}
	if stddev != 2.0 {
		t.Errorf("got=%f want=2.0", stddev)
	}
}
