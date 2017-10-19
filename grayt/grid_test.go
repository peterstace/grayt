package grayt

import (
	"testing"
)

func TestGridPopulationCrash(t *testing.T) {
	// Causesed a crash before bugfix.
	objs := Group(
		AlignedSquare(Vect(0, 1, 0), Vect(1, 1, -1)),
		AlignedSquare(Vect(-10, -10, 1.3), Vect(10, 10, 1.3)),
	)
	newGrid(4, objs)
}
