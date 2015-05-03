package grayt

type Ray struct {
	Start, Dir Vect
}

func (r Ray) At(t float64) Vect {
	return r.Start.Add(r.Dir.Extended(t))
}
