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
		Cam: grayt.NewRectilinearCamera(),
	}
}

func main() {

	quality := grayt.DefaultQuality()

	err := grayt.TraceAnimation(sceneFactory{}, "output", &quality)
	if err != nil {
		log.Fatal(err)
	}
}
