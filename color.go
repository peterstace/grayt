package grayt

import (
	"image/color"
	"math"
)

func float64ToUint16(x float64) uint16 {
	if x < 0.0 {
		return 0
	}
	if x > 1.0 {
		return math.MaxUint16
	}
	return uint16(x * float64(math.MaxUint16))
}

func NewColor(r, g, b float64) color.NRGBA64 {
	return color.NRGBA64{
		R: float64ToUint16(r),
		G: float64ToUint16(g),
		B: float64ToUint16(b),
		A: math.MaxUint16,
	}
}

// TODO: Scale and Multiply
