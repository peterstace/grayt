package main

import "github.com/peterstace/grayt"

func main() {

	tracer := grayt.NewTracer()
	sceneFactory := func(t float64) grayt.Scene {
		return grayt.Scene{}
	}
	tracer.Trace("output", sceneFactory)
}
