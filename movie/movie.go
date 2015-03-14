package movie

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"math/rand"
	"os"
	"path"

	"github.com/peterstace/grayt/tracer"
)

type Camera func(float64) tracer.Camera
type Geometry func(float64) tracer.Geometry
type Light func(float64) tracer.Light

type Movie struct {
	Frames     int
	Camera     Camera
	Geometries []Geometry
	Lights     []Light
}

func TraceMovie(m Movie, filename string) error {

	// Get a temp dir to put the JPEG for each frame into.
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	for i := 0; i < m.Frames; i++ {

		// Create the sample(s).
		t := Sample(i, m.Frames)
		var sample tracer.Scene
		sample.Camera = m.Camera(t)
		for _, g := range m.Geometries {
			sample.Geometries = append(sample.Geometries, g(t))
		}
		for _, l := range m.Lights {
			sample.Lights = append(sample.Lights, l(t))
		}

		// Trace the image.
		img := tracer.TraceImage([]tracer.Scene{sample})

		// Write the image out to file.
		if out, err := os.Create(path.Join(tmpDir, fmt.Sprintf("%d.jpg", i))); err != nil {
			return err
		} else if err := jpeg.Encode(out, img, nil); err != nil {
			return err
		}
	}

	// Create a movie file from the images.
	// XXX

	return nil
}

// Sample returns a value in the interval [0, 1). The interval [0, 1) is
// divided equally into n parts, with the random value being selected from
// within the ith segment. Precondition: 0 <= i < n.
func Sample(i int, n int) float64 {
	return (float64(i) + rand.Float64()) / float64(n)
}
