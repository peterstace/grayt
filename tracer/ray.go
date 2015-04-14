package tracer

type Ray struct {
	Start, Dir Vect
}

func (r Ray) At(t float64) Vect {
	return VectAdd(r.Start, r.Dir.Extended(t))
}
