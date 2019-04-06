package trace

import (
	"encoding/binary"
	"image"
	"io"
	"sync"

	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/xmath"
)

type accumulator struct {
	mu        sync.Mutex
	passes    int
	dim       xmath.Dimensions
	aggregate []colour.Colour
	landing   []colour.Colour
}

func newAccumulator(dim xmath.Dimensions) *accumulator {
	acc := new(accumulator)
	acc.dim = dim
	n := dim.Wide * dim.High
	acc.aggregate = make([]colour.Colour, n)
	acc.landing = make([]colour.Colour, n)
	return acc
}

func (a *accumulator) set(x, y int, c colour.Colour) {
	idx := x + a.dim.Wide*y
	a.landing[idx] = c
}

func (a *accumulator) merge(depth int) {
	a.mu.Lock()
	a.passes += depth
	for i, c := range a.landing {
		a.aggregate[i] = a.aggregate[i].Add(c)
	}
	a.mu.Unlock()
}

func (a *accumulator) getPasses() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.passes
}

// toImage converts the accumulator into an image. Exposure controls how bright
// the arithmetic mean brightness in the image is. A value of 1.0 results in a
// mean brightness half way between black and white.
func (a *accumulator) toImage(exposure float64) image.Image {
	a.mu.Lock()
	const gamma = 2.2
	mean := a.mean()
	img := image.NewNRGBA(image.Rect(0, 0, a.dim.Wide, a.dim.High))
	for x := 0; x < a.dim.Wide; x++ {
		for y := 0; y < a.dim.High; y++ {
			i := y*a.dim.Wide + x
			img.Set(x, y, a.aggregate[i].
				Scale(0.5*exposure/mean).
				Pow(1.0/gamma).
				ToNRGBA())
		}
	}
	a.mu.Unlock()
	return img
}

func (a *accumulator) mean() float64 {
	var sum float64
	for _, c := range a.aggregate {
		sum += c.R + c.G + c.B
	}
	return sum / float64(len(a.aggregate)) / 3.0
}

type countingWriter struct {
	w io.Writer
	n int
}

func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.n += n
	return n, err
}

type countingReader struct {
	r io.Reader
	n int
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += n
	return n, err
}

func (a *accumulator) WriteTo(w io.Writer) (int64, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cw := countingWriter{w, 0}
	for _, data := range []interface{}{
		int64(a.dim.Wide), int64(a.dim.High), int64(a.passes), a.aggregate,
	} {
		if err := binary.Write(&cw, binary.BigEndian, data); err != nil {
			return int64(cw.n), err
		}
	}
	return int64(cw.n), nil
}

func (a *accumulator) ReadFrom(r io.Reader) (int64, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cr := countingReader{r, 0}
	for _, data := range []interface{}{&a.dim, &a.passes, &a.aggregate} {
		if err := binary.Read(&cr, binary.BigEndian, data); err != nil {
			return int64(cr.n), err
		}
	}
	return int64(cr.n), nil
}
