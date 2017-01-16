package grayt

import "sync"

type accumulator struct {
	pixels []pixel
	wide   int
	high   int
}

func newAccumulator(wide, high int) *accumulator {
	return &accumulator{
		pixels: make([]pixel, wide*high),
		wide:   wide,
		high:   high,
	}
}

func (a *accumulator) add(x, y int, c Colour, index int) {
	i := y*a.wide + x
	a.pixels[i].mu.Lock()
	a.pixels[i].add(c, index)
	a.pixels[i].mu.Unlock()
}

func (a *accumulator) get(x, y int) Colour {
	i := y*a.wide + x
	return a.pixels[i].colourSum
}

type pixel struct {
	mu        sync.Mutex
	colourSum Colour
	nextIndex int
	pending   []pendingColour
}

type pendingColour struct {
	colour Colour
	index  int
}

func (e *pixel) add(c Colour, index int) {

	if e.nextIndex != index {
		e.pending = append(e.pending, pendingColour{c, index})
		return
	}

	e.colourSum = e.colourSum.add(c)
	e.nextIndex++

	for true {
		found := false
		for i := range e.pending {
			if e.pending[i].index == e.nextIndex {
				found = true
				e.colourSum = e.colourSum.add(e.pending[i].colour)
				e.nextIndex++
				e.pending[i] = e.pending[len(e.pending)-1]
				e.pending = e.pending[:len(e.pending)-1]
				break
			}
		}
		if !found {
			break
		}
	}
}
