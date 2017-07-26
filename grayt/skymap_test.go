package grayt

import "testing"

func TestAzimuth(t *testing.T) {
	for _, test := range []struct {
		dir     Vector
		azimuth float64
	}{
		{Vect(0, 0, -1), 0.5},
	} {
		actual := azimuth(test.dir)
		if actual != test.azimuth {
			t.Errorf("for %v: want=%v got=%v", test.dir, test.azimuth, actual)
		}
	}
}
