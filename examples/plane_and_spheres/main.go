package main

import (
	"log"

	"github.com/peterstace/grayt"
)

type sceneFactory struct{}

func (f sceneFactory) FrameCount() int {
	return 1
}

func (f sceneFactory) MakeScene(t float64) grayt.Scene {
	return grayt.Scene{
		Camera: grayt.NewRectilinearCamera(),
	}
}

func main() {

	config := grayt.DefaultConfig()
	tracer := grayt.NewAnimationTracer(config)

	err := tracer.TraceAnimation("output", sceneFactory{})
	if err != nil {
		log.Fatal(err)
	}
}
