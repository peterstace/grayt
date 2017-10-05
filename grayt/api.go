package grayt

type API struct {
	aspectRatio vect2
}

// PushMatrix duplicates the matrix at the top of the stack and pushes it.
func (a *API) PushMatrix() {
	// TODO
}

// PopMatrix pops a matrix from the top of the stack. If the stack becomes
// empty, it panics.
func (a *API) PopMatrix() {
	// TODO
}

// Translate translates the matrix at the top of the stack.
func (a *API) Translate(x, y, z float64) {
	// TODO
}
