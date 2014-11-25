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
	numFrames       int
	samplesPerFrame int
	engine          Engine
}

func NewAnimationTracer() AnimationBuilder {
	return AnimationBuilder{
		numFrames:       1,
		samplesPerFrame: 1,
		engine: Engine{
			pxWide: 640,
			pxHigh: 480,
		},
	}
}

func (b *AnimationBuilder) SetNumFrames(numFrames int) {
	b.numFrames = numFrames
}

func (b *AnimationBuilder) SetSamplesPerFrame(samplesPerFrame int) {
	b.samplesPerFrame = samplesPerFrame
}

func (b *AnimationBuilder) GetEngine() *Engine {
	return &b.engine
}

func (b *AnimationBuilder) TraceAnimation(
	outDir string,
	sceneFactory func(float64) Scene,
) error {

	log.Print("Tracing Animation...")
	animationStartTime := time.Now()

	for i := 0; i < b.numFrames; i++ {
		frameStartTime := time.Now()

		// Create scenes
		scenes := make([]Scene, b.samplesPerFrame)
		for j := 0; j < b.samplesPerFrame; j++ {
			scenes[j] = sceneFactory(b.calculateSampleOffset(i, j))
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
			i+1, b.numFrames, time.Now().Sub(frameStartTime))
	}

	log.Printf("Done (%v)", time.Now().Sub(animationStartTime))

	return nil
}

func (b *AnimationBuilder) calculateSampleOffset(frame, sample int) float64 {
	frameWidth := 1.0 / float64(b.numFrames)
	sampleWidth := frameWidth / float64(b.samplesPerFrame)
	return float64(frame)*frameWidth + (float64(sample)+rand.Float64())*sampleWidth
}
