package grayt

import (
	"fmt"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

// TraceAnimation produces a series of JPG images using scenes produced using a
// SceneFactory.
func TraceAnimation(
	sceneFactory SceneFactory,
	outDir string,
	quality *Quality,
) error {

	log.Print("Tracing Animation...")
	animationStartTime := time.Now()

	for i := 0; i < sceneFactory.FrameCount(); i++ {
		frameStartTime := time.Now()

		// Create samples
		samples := make([]Scene, quality.TemporalAALevel)
		for j := range samples {
			offset := calculateSampleOffset(i, j,
				sceneFactory.FrameCount(),
				quality.TemporalAALevel)
			samples[j] = sceneFactory.MakeScene(offset)
		}

		// Trace the samples.
		img := TraceScene(quality, samples...)

		// Write to file.
		filename := path.Join(outDir, fmt.Sprintf("%d.jpeg", i))
		file, err := os.Create(filename)
		if err != nil {
			file.Close()
			return err
		}
		err = jpeg.Encode(file, img, nil)
		if err != nil {
			file.Close()
			return err
		}
		file.Close()

		// Log out that we're done!
		log.Printf("Frame %d of %d complete (%v)",
			i+1, sceneFactory.FrameCount(), time.Now().Sub(frameStartTime))
	}

	log.Printf("Done (%v)", time.Now().Sub(animationStartTime))

	return nil
}

func calculateSampleOffset(frame, sample, frameCount, temporalAALevel int) float64 {
	frameWidth := 1.0 / float64(frameCount)
	sampleWidth := frameWidth / float64(temporalAALevel)
	return float64(frame)*frameWidth + (float64(sample)+rand.Float64())*sampleWidth
}
