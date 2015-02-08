package ray

import "github.com/peterstace/grayt/vect"

type Ray struct {
	Start, Dir vect.Vect
}

func (r Ray) At(t float64) vect.Vect {
	return vect.Add(r.Start, r.Dir.Extended(t))
}
