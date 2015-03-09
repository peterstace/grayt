package main

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/peterstace/grayt/tracer"
	"github.com/peterstace/grayt/vect"
)

func main() {

	scene := tracer.Scene{
		Camera: tracer.NewRectilinearCamera(
			tracer.CameraConfig{
				Location:      vect.New(0, 0, 0),
				ViewDirection: vect.New(0, 0, -1),
				UpDirection:   vect.New(0, 1, 0),
				FieldOfView:   1.5,
				FocalLength:   10.0,
				FocalRatio:    1000.0,
			}),
		Geometries: []tracer.Geometry{
			tracer.NewPlane(vect.New(0, -2, 0), vect.New(0, 1, 0)),
			tracer.NewSphere(vect.New(0, 0, -10), 1),
		},
		Lights: []tracer.Light{
			tracer.Light{Location: vect.New(0, 10, -10), Radius: 0.5, Intensity: 100},
		},
	}

	img := tracer.TraceImage([]tracer.Scene{scene})

	if out, err := os.Create("out.jpg"); err != nil {
		log.Fatal(err)
	} else if err := jpeg.Encode(out, img, nil); err != nil {
		log.Fatal(err)
	}
}
