package grayt

type Vect struct {
	X, Y, Z float64
}

type ray struct {
	start, dir Vect
}
