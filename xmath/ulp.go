package xmath

import "math"

func float64AsInt64(f float64) int64 {
	if f < 0 {
		return -int64(math.Float64bits(-f))
	} else {
		return int64(math.Float64bits(f))
	}
}

func int64AsFloat64(i int64) float64 {
	if i < 0 {
		return -math.Float64frombits(uint64(-i))
	} else {
		return math.Float64frombits(uint64(i))
	}
}

func AddULPs(f float64, ulps int64) float64 {
	return int64AsFloat64(float64AsInt64(f) + ulps)
}
