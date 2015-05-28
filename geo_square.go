package grayt

func NewSquare(m Material, v1, v2, v3, v4 Vect) []Geometry {
	return []Geometry{
		NewTriangle(m, v1, v2, v3),
		NewTriangle(m, v3, v4, v1),
	}
}
