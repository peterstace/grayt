package trace

import (
	"testing"
)

func TestGridPopulationCrash(t *testing.T) {
	//Causesed a crash before bugfix.
	objs := []object{
		{
			Surface: &alignYSquare{
				X1: 0,
				X2: 1,
				Y:  1,
				Z1: 0,
				Z2: -1,
			},
		},
		{
			Surface: &alignZSquare{
				X1: -10,
				X2: 10,
				Y1: -10,
				Y2: 10,
				Z:  1.3,
			},
		},
	}
	newGrid(4, objs)
}
