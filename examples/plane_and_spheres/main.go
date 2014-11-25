package main

import (
	"log"

	"github.com/peterstace/grayt"
)

func main() {

	tracer := grayt.NewAnimationTracer()
	tracer.SetNumFrames(10)
	tracer.SetSamplesPerFrame(3)

	sceneFactory := func(t float64) grayt.Scene {
		return grayt.Scene{}
	}
	err := tracer.TraceAnimation("output", sceneFactory)
	if err != nil {
		log.Fatal(err)
	}
}
