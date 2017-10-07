package engine

func triangle(e, d vect3) (n, h vect3, ok bool) {
	const (
		invSqrt3 = 1.0 / 1.73205080756887729352744634150587236694280525381038
		third    = 1.0 / 3.0
	)
	n = vect3{invSqrt3, invSqrt3, invSqrt3}
	t := vect3{third, third, third}.sub(e).dot(n) / d.dot(n)
	if t < 0 {
		return
	}
	h = e.add(d.scale(t))
	if h[0] < 0 || h[1] < 0 || h[2] < 0 || h[0]+h[1]+h[2] < 1 {
		return
	}
	ok = true
	return
}
