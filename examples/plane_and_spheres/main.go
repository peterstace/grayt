package main

import "github.com/peterstace/grayt"

func main() {
	sceneFactory := func(t float64) grayt.Scene {
		return grayt.Scene{}
	}
	grayt.TraceAnimation("output", sceneFactory)
}
