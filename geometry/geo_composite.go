package grayt

type Square struct {
	V1, V2, V3, V4 Vector
}

func (s Square) MakeSurfaces() []Surface {
	return append(
		Triangle{s.V1, s.V2, s.V3}.MakeSurfaces(),
		Triangle{s.V3, s.V4, s.V1}.MakeSurfaces()...,
	)
}
