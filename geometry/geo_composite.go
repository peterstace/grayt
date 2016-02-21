package geometry

import "github.com/peterstace/grayt"

type Square struct {
	V1, V2, V3, V4 grayt.Vector
}

func (s Square) MakeSurfaces() []grayt.Surface {
	return append(
		Triangle{s.V1, s.V2, s.V3}.MakeSurfaces(),
		Triangle{s.V3, s.V4, s.V1}.MakeSurfaces()...,
	)
}
