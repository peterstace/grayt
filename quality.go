package grayt

type Quality struct {
	pxWide, pxHigh  int
	temporalAALevel int
	spatialAALevel  int
}

func (q *Quality) PxWide() int {
	if q == nil {
		return 320
	}
	return q.pxWide
}

func (q *Quality) PxHigh() int {
	if q == nil {
		return 240
	}
	return q.pxHigh
}

func (q *Quality) TemporalAALevel() int {
	if q == nil {
		return 1
	}
	return q.temporalAALevel
}

func (q *Quality) SpatialAALevel() int {
	if q == nil {
		return 1
	}
	return q.spatialAALevel
}
