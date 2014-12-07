package grayt

type Quality struct {
	PxWide, PxHigh  int
	TemporalAALevel int
	SpatialAALevel  int
}

func DefaultQuality() Quality {
	return Quality{
		PxWide:          320,
		PxHigh:          240,
		TemporalAALevel: 1,
		SpatialAALevel:  1,
	}
}
