package ray

import "github.com/peterstace/grayt/vect"

type Ray struct {
	Start, Dir vect.V
}

func (r Ray) At(t float64) vect.V {
	return vect.Add(r.Start, r.Dir.Extended(t))
}
