package tracer

import (
	"fmt"
	"image/jpeg"
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
		img := TraceImage([]Scene{sample})

		// Output the image.
		// TODO: Should pad filename with leading zeros.
		if file, err := os.Create(path.Join(outputDir, fmt.Sprintf("%d.jpg", i))); err != nil {
			return err
		} else if err := jpeg.Encode(file, img, nil); err != nil {
			return err
		}
	}
	return nil
}

//func outputMovie(imgs []image.Image, fps int, filename string) error {
//
//	log.Print("outputting movie")
//
//	// Get a temp dir to put the JPEG for each frame into.
//	tmpDir, err := ioutil.TempDir("", "")
//	if err != nil {
//		return err
//	}
//	defer func() {
//		log.Printf("cleaning up %q", tmpDir)
//		os.RemoveAll(tmpDir)
//	}()
//
//	// Write the images out to file.
//	for i := 0; i < len(imgs); i++ {
//		if out, err := os.Create(path.Join(tmpDir, fmt.Sprintf("%d.jpg", i))); err != nil {
//			return err
//		} else if err := jpeg.Encode(out, imgs[i], nil); err != nil {
//			return err
//		}
//	}
//
//	// Create a movie file from the images.
//	wd, err := os.Getwd()
//	if err != nil {
//		return err
//	}
//	cmd := exec.Command("ffmpeg",
//		"-framerate", fmt.Sprintf("%d", fps),
//		"-i", "%d.jpg",
//		"-codec", "copy",
//		path.Join(wd, filename))
//	cmd.Dir = tmpDir
//	out, err := cmd.CombinedOutput()
//	if err != nil {
//		return fmt.Errorf("%s: %s", err, string(out))
//	}
//
//	log.Printf("movie output to %q", filename)
//	return nil
//}

// Sample returns a value in the interval [0, 1). The interval [0, 1) is
// divided equally into n parts, with the random value being selected from
// within the ith segment. Precondition: 0 <= i < n.
func Sample(i int, n int) float64 {
	return (float64(i) + rand.Float64()) / float64(n)
}
