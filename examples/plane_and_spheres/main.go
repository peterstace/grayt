package main

import (
	"log"

	"github.com/peterstace/grayt"
)

func main() {

	config := grayt.DefaultConfig()
	tracer := grayt.NewAnimationTracer(config)

	sceneFactory := func(t float64) grayt.Scene {
		return grayt.Scene{}
	}
	err := tracer.TraceAnimation("output", sceneFactory)
	if err != nil {
		log.Fatal(err)
	}
}
