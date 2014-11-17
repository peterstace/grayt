package main

import "github.com/peterstace/grayt"

func main() {

	tracer := grayt.NewRayTracer()
	tracer.SetNumFrames(10)
	tracer.SetSamplesPerFrame(3)

	sceneFactory := func(t float64) grayt.Scene {
		return grayt.Scene{}
	}
	tracer.TraceAnimation("output", sceneFactory)
}
