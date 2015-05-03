package grayt

import "unsafe"

func ulpDiff(a, b float64) uint64 {

	ulpA := float64ToULP(a)
	ulpB := float64ToULP(b)

	if ulpA > ulpB {
		return ulpA - ulpB
	} else {
		return ulpB - ulpA
	}
}

func float64ToULP(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}
