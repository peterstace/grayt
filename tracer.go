package grayt

import "log"

type Tracer struct {
	numFrames       int
	samplesPerFrame int
}

func NewTracer() Tracer {
	return Tracer{
		numFrames:       1,
		samplesPerFrame: 1,
	}
}

func (t *Tracer) SetSamplesPerFrame(samplesPerFrame int) {
	t.samplesPerFrame = samplesPerFrame
}

func (t *Tracer) Trace(path string, sceneFactory func(t float64) Scene) error {

	log.Print("Tracing Animation...")

	return nil
}
