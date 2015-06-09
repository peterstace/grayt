package grayt

func NewSquare(v1, v2, v3, v4 Vect) []Surface {
	return []Surface{
		NewTriangle(v1, v2, v3),
		NewTriangle(v3, v4, v1),
	}
}
