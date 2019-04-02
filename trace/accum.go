package trace

import (
	"image"
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

/*
func loadAccumulator(fname string) (string, *accumulator, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	var (
		sceneNameLen int64
		passes       int64
		wide         int64
		high         int64
	)
	if err := binary.Read(f, binary.LittleEndian, &sceneNameLen); err != nil {
		return "", nil, err
	}
	sceneName := make([]byte, sceneNameLen)
	if err := binary.Read(f, binary.LittleEndian, &sceneName); err != nil {
		return "", nil, err
	}
	if err := binary.Read(f, binary.LittleEndian, &passes); err != nil {
		return "", nil, err
	}
	if err := binary.Read(f, binary.LittleEndian, &wide); err != nil {
		return "", nil, err
	}
	if err := binary.Read(f, binary.LittleEndian, &high); err != nil {
		return "", nil, err
	}
	pixels := make([]colour.Colour, wide*high)
	if err := binary.Read(f, binary.LittleEndian, &pixels); err != nil {
		return "", nil, err
	}

	return string(sceneName), &accumulator{
		passes: passes,
		pixelGrid: pixelGrid{
			wide:   int(wide),
			high:   int(high),
			pixels: pixels,
		},
	}, nil
}
*/

/*
func (a *accumulator) save(fname, sceneName string) error {
	a.Lock()
	defer a.Unlock()

	f, err := ioutil.TempFile(".", "")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	if err := binary.Write(f, binary.LittleEndian, int64(len(sceneName))); err != nil {
		f.Close()
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, []byte(sceneName)); err != nil {
		f.Close()
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, int64(a.passes)); err != nil {
		f.Close()
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, int64(a.wide)); err != nil {
		f.Close()
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, int64(a.high)); err != nil {
		f.Close()
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, a.pixels); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), fname)
}
*/
