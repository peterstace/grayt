package main

import (
	"image"
	"log"
	"os"

	_ "image/jpeg"

	. "github.com/peterstace/grayt/grayt"
)

func scene() Scene {
	loc := Vect(-2, 1, 4)
	lookAt := Vect(0, 0, 0)
	return Scene{
		Camera: Camera().With(
			Location(loc),
			LookingAt(lookAt),
			FieldOfViewInDegrees(90),
			AspectRatioWidthAndHeight(4, 3),
			FocalLengthAndRatio(lookAt.Sub(loc).Length(), 25),
		),
		Objects: Group(
			AlignedBox(Vect(-1, -5, -1), Vect(1, -1, 1)),
			Sphere(Vect(0, -0.2, 0), 0.8).With(Mirror()),
		),
	}
}

func main() {
	scn := scene()
	scn.Sky = MustCreateSkymap("/media/sf_shared/HDRI_cotty_hennestrand_equirectangular_690x345.jpg")
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
