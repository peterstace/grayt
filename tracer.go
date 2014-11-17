package grayt

import (
	"log"
	"math/rand"
	"time"
)

type RayTracer struct {
	numFrames       int
	samplesPerFrame int
}

func NewRayTracer() RayTracer {
	return RayTracer{
		numFrames:       1,
		samplesPerFrame: 1,
	}
}

func (r *RayTracer) SetNumFrames(numFrames int) {
	r.numFrames = numFrames
}

func (r *RayTracer) SetSamplesPerFrame(samplesPerFrame int) {
	r.samplesPerFrame = samplesPerFrame
}

func (r *RayTracer) TraceAnimation(path string, sceneFactory func(float64) Scene) error {

	log.Print("Tracing Animation...")
	animationStartTime := time.Now()

	for i := 0; i < r.numFrames; i++ {
		frameStartTime := time.Now()

		// Create scenes
		scenes := make([]Scene, r.samplesPerFrame)
		for j := 0; j < r.samplesPerFrame; j++ {
			scenes[j] = sceneFactory(r.calculateT(i, j))
		}

		log.Printf("Frame %d of %d complete (%v)",
			i+1, r.numFrames, time.Now().Sub(frameStartTime))
	}

	log.Printf("Done (%v)", time.Now().Sub(animationStartTime))

	return nil
}

func (r *RayTracer) calculateT(frame, sample int) float64 {
	frameWidth := 1.0 / float64(r.numFrames)
	sampleWidth := frameWidth / float64(r.samplesPerFrame)
	return float64(frame)*frameWidth + (float64(sample)+rand.Float64())*sampleWidth
}
