package main

import (
	"image"
	"log"
	"os"

	. "github.com/peterstace/grayt/grayt"
)

func scene() Scene {
	return Scene{
		Camera: Camera().With(
			Location(Vect(0, 0, 0)),
			LookingAt(Vect(0, 0, -1)),
			FieldOfViewInDegrees(90),
			AspectRatioWidthAndHeight(4, 3),
		),
		Objects: Group(),
	}
}

func main() {
	scn := scene()
	scn.Sky = MustCreateSkymap("/media/sf_shared/custom_skymap.png")
	Run("clouds", scn)
}

func MustCreateSkymap(filename string) Skymap {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	sky, err := CreateSkymap(img)
	if err != nil {
		log.Fatal(err)
	}
	return sky
}
