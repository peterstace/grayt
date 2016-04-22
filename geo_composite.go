package grayt

func Square(v1, v2, v3, v4 Vector) []Surface {
	return []Surface{
		Triangle(v1, v2, v3),
		Triangle(v3, v4, v1),
	}
}
