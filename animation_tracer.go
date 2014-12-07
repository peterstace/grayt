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

type AnimationBuilder struct {
	engine  engine
	quality Quality
}

func NewAnimationTracer(q Quality) AnimationBuilder {
	return AnimationBuilder{
		engine:  engine{quality: q},
		quality: q,
	}
}

func (b *AnimationBuilder) TraceAnimation(outDir string, sceneFactory SceneFactory) error {

	log.Print("Tracing Animation...")
	animationStartTime := time.Now()

	for i := 0; i < sceneFactory.FrameCount(); i++ {
		frameStartTime := time.Now()

		// Create scenes
		scenes := make([]Scene, b.quality.TemporalAALevel)
		for j := range scenes {
			offset := b.calculateSampleOffset(i, j, sceneFactory.FrameCount())
			scenes[j] = sceneFactory.MakeScene(offset)
		}

		// Trace the scenes.
		img := b.engine.traceScenes(scenes)

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

func (b *AnimationBuilder) calculateSampleOffset(frame, sample, frameCount int) float64 {
	frameWidth := 1.0 / float64(frameCount)
	sampleWidth := frameWidth / float64(b.quality.TemporalAALevel)
	return float64(frame)*frameWidth + (float64(sample)+rand.Float64())*sampleWidth
}
