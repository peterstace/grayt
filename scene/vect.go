package scene

type Vector struct {
	X, Y, Z float64
}

func Vect(x, y, z float64) Vector {
	return Vector{x, y, z}
}
