package geometry

import "github.com/peterstace/grayt"

func Square(v1, v2, v3, v4 grayt.Vector) []grayt.Surface {
	return []grayt.Surface{
		Triangle(v1, v2, v3),
		Triangle(v3, v4, v1),
	}
}
