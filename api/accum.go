package api

import (
	"encoding/binary"
	"image"
	"io/ioutil"
	"os"
	"sync"

	"github.com/peterstace/grayt/colour"
)

// TODO: this file needs a big review

type pixelGrid struct {
	wide   int
	high   int
	pixels []colour.Colour
}

func (g *pixelGrid) set(x, y int, c colour.Colour) {
	i := y*g.wide + x
	g.pixels[i] = c
}

type accumulator struct {
	sync.Mutex
	passes int64
	pixelGrid
}

func newAccumulator(pxWide, pxHigh int) *accumulator {
	acc := new(accumulator)
	acc.wide = pxWide
	acc.high = pxHigh
	acc.pixels = make([]colour.Colour, pxWide*pxHigh)
	return acc
}

func (a *accumulator) getPasses() int64 {
	a.Lock()
	p := a.passes
	a.Unlock()
	return p
}

func (a *accumulator) merge(g *pixelGrid) {
	a.Lock()
	defer a.Unlock()

	a.passes++
	for i, c := range a.pixels {
		a.pixels[i] = c.Add(g.pixels[i])
	}
}

// toImage converts the accumulator into an image. Exposure controls how bright
// the arithmetic mean brightness in the image is. A value of 1.0 results in a
// mean brightness half way between black and white.
func (a *accumulator) toImage(exposure float64) image.Image {
	a.Lock()
	defer a.Unlock()

	const gamma = 2.2
	mean := a.mean()
	img := image.NewNRGBA(image.Rect(0, 0, a.wide, a.high))
	for x := 0; x < a.wide; x++ {
		for y := 0; y < a.high; y++ {
			i := y*a.wide + x
			img.Set(x, y, a.pixels[i].
				Scale(0.5*exposure/mean).
				Pow(1.0/gamma).
				ToNRGBA())
		}
	}
	return img
}

func (a *accumulator) mean() float64 {
	var sum float64
	for _, c := range a.pixels {
		sum += c.R + c.G + c.B
	}
	return sum / float64(len(a.pixels)) / 3.0
}

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
