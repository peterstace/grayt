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
		Camera: grayt.NewRectilinearCamera(grayt.CameraConfig{

			Location:      grayt.Vect{},
			ViewDirection: grayt.Vect{0.0, 0.0, -1.0},
			UpDirection:   grayt.Vect{0.0, 1.0, 0.0},
			FieldOfView:   90,
			FocalLength:   1.0,
			FocalRatio:    1.0,
		}),
	}
}

func main() {

	quality := grayt.DefaultQuality()

	err := grayt.TraceAnimation(sceneFactory{}, "output", &quality)
	if err != nil {
		log.Fatal(err)
	}
}
