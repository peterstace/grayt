package grayt

import "math"

func ulpDiff(a, b float64) uint64 {

	ulpA := math.Float64bits(a)
	ulpB := math.Float64bits(b)

	if ulpA > ulpB {
		return ulpA - ulpB
	} else {
		return ulpB - ulpA
	}
}
