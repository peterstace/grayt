package tracer

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"time"
)

type CameraFn func(float64) Camera
type GeometryFn func(float64) Geometry
type LightFn func(float64) Light

func StillCamera(c Camera) CameraFn {
	return func(float64) Camera { return c }
}
func StillGeometry(g Geometry) GeometryFn {
	return func(float64) Geometry { return g }
}
func StillLight(l Light) LightFn {
	return func(float64) Light { return l }
}

type Movie struct {
	Duration   time.Duration
	Camera     CameraFn
	Geometries []GeometryFn
	Lights     []LightFn
}

func TraceMovie(m Movie, filename string) error {

	// Calculate the number of frames.
	const fps = 24
	numFrames := int(m.Duration.Seconds() * fps)

	imgs := make([]image.Image, numFrames)
	for i := 0; i < numFrames; i++ {

		// Create the sample(s).
		t := Sample(i, numFrames)
		var sample Scene
		sample.Camera = m.Camera(t)
		for _, g := range m.Geometries {
			sample.Geometries = append(sample.Geometries, g(t))
		}
		for _, l := range m.Lights {
			sample.Lights = append(sample.Lights, l(t))
		}

		// Trace the image.
		imgs[i] = TraceImage([]Scene{sample})

	}

	return outputMovie(imgs, fps, filename)
}

func outputMovie(imgs []image.Image, fps int, filename string) error {

	log.Print("outputting movie")

	// Get a temp dir to put the JPEG for each frame into.
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer func() {
		log.Printf("cleaning up %q", tmpDir)
		os.RemoveAll(tmpDir)
	}()

	// Write the images out to file.
	for i := 0; i < len(imgs); i++ {
		if out, err := os.Create(path.Join(tmpDir, fmt.Sprintf("%d.jpg", i))); err != nil {
			return err
		} else if err := jpeg.Encode(out, imgs[i], nil); err != nil {
			return err
		}
	}

	// Create a movie file from the images.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd := exec.Command("ffmpeg",
		"-framerate", fmt.Sprintf("%d", fps),
		"-i", "%d.jpg",
		"-codec", "copy",
		path.Join(wd, filename))
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}

	log.Printf("movie output to %q", filename)
	return nil
}

// Sample returns a value in the interval [0, 1). The interval [0, 1) is
// divided equally into n parts, with the random value being selected from
// within the ith segment. Precondition: 0 <= i < n.
func Sample(i int, n int) float64 {
	return (float64(i) + rand.Float64()) / float64(n)
}
