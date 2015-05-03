package grayt

import (
	"fmt"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"path"
)

type Movie struct {
	FrameCount        int
	CameraFactory     func(t float64) Camera
	GeometriesFactory func(t float64) []Geometry
	LightsFactory     func(t float64) []Light
}

func TraceMovie(m Movie, outputDir string) error {

	for i := 0; i < m.FrameCount; i++ {

		// Create the sample(s).
		// TODO: This currently only creates a single sample per frame. It
		// should instead generate multiple samples per frame.
		t := Sample(i, m.FrameCount)
		sample := Scene{
			Camera:     m.CameraFactory(t),
			Geometries: m.GeometriesFactory(t),
			Lights:     m.LightsFactory(t),
		}

		// Trace the image.
		log.Printf("Tracing image %d of %d", i, m.FrameCount)
		img := TraceImage([]Scene{sample})

		// Output the image.
		// TODO: Should pad filename with leading zeros.
		numWidth := len(fmt.Sprintf("%d", m.FrameCount))
		filepath := path.Join(outputDir, fmt.Sprintf("%0*d.jpg", numWidth, i))
		log.Printf("Saving traced image to %q", filepath)
		if file, err := os.Create(filepath); err != nil {
			return err
		} else if err := jpeg.Encode(file, img, nil); err != nil {
			return err
		}
	}
	return nil
}

// Sample returns a value in the interval [0, 1). The interval [0, 1) is
// divided equally into n parts, with the random value being selected from
// within the ith segment. Precondition: 0 <= i < n.
func Sample(i int, n int) float64 {
	return (float64(i) + rand.Float64()) / float64(n)
}
