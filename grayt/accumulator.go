package grayt

import (
	"encoding/binary"
	"fmt"
	"image"
	"io/ioutil"
	"os"
)

type pixelGrid struct {
	wide   int
	high   int
	pixels []Colour
}

func (g *pixelGrid) set(x, y int, c Colour) {
	i := y*g.wide + x
	g.pixels[i] = c
}

type accumulator struct {
	// TODO: Store hash as well so scenes don't get mixed up.
	count int
	pixelGrid
}

func (a *accumulator) merge(g *pixelGrid) {
	a.count++
	for i, c := range a.pixels {
		a.pixels[i] = c.add(g.pixels[i])
	}
}

func (a *accumulator) mean() float64 {
	var sum float64
	for _, c := range a.pixels {
		sum += c.R + c.G + c.B
	}
	return sum / float64(len(a.pixels)) / 3.0
}

// ToImage converts the accumulator into an image. Exposure controls how bright
// the arithmetic mean brightness in the image is. A value of 1.0 results in a
// mean brightness half way between black and white.
func (a *accumulator) toImage(exposure float64) image.Image {
	const gamma = 2.2
	mean := a.mean()
	img := image.NewNRGBA(image.Rect(0, 0, a.wide, a.high))
	for x := 0; x < a.wide; x++ {
		for y := 0; y < a.high; y++ {
			i := y*a.wide + x
			img.Set(x, y, a.pixels[i].
				scale(0.5*exposure/mean).
				pow(1.0/gamma).
				toNRGBA())
		}
	}
	return img
}

func (a *accumulator) load() (bool, error) {
	f, err := os.Open("checkpoint")
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return false, nil
		}
		return false, err
	}
	defer f.Close()

	var count int64
	if err := binary.Read(f, binary.LittleEndian, &count); err != nil {
		return false, err
	}
	a.count = int(count)

	var wide int64
	if err := binary.Read(f, binary.LittleEndian, &wide); err != nil {
		return false, err
	}
	a.wide = int(wide)

	var high int64
	if err := binary.Read(f, binary.LittleEndian, &high); err != nil {
		return false, err
	}
	a.high = int(high)

	fmt.Println(wide, high)
	a.pixels = make([]Colour, wide*high)
	if err := binary.Read(f, binary.LittleEndian, &a.pixels); err != nil {
		return false, err
	}
	return true, nil
}

func (a *accumulator) save() error {
	f, err := ioutil.TempFile(".", "")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err := binary.Write(f, binary.LittleEndian, int64(a.count)); err != nil {
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
	return os.Rename(f.Name(), "checkpoint")
}
