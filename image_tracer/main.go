package main

import (
	"log"
	"time"

	"github.com/peterstace/grayt/tracer"
	"github.com/peterstace/grayt/vect"
)

func main() {

	mov := tracer.Movie{
		Duration: time.Second * 17,
		Camera: tracer.StillCamera(tracer.NewRectilinearCamera(
			tracer.CameraConfig{
				Location:      vect.New(0, 10, 5),
				ViewDirection: vect.New(0, -1, -2),
				UpDirection:   vect.New(0, 1, 0),
				FieldOfView:   1.5,
				FocalLength:   10.0,
				FocalRatio:    1000.0,
			})),
		Geometries: []tracer.GeometryFn{

			tracer.StillGeometry(tracer.NewPlane(vect.New(0, 0, 0), vect.New(0, 1, 0))),

			tracer.StillGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(0.2, 1, 0))),
			tracer.StillGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(-0.2, 1, 0))),
			tracer.StillGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(0, 1, 0.2))),
			tracer.StillGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(0, 1, -0.2))),

			tracer.StillGeometry(tracer.NewSphere(vect.New(0, 1, -10), 1)),
		},

		Lights: []tracer.LightFn{
			tracer.StillLight(tracer.Light{Location: vect.New(0, 10, -10), Radius: 0.5, Intensity: 100}),
		},
	}

	if err := tracer.TraceMovie(mov, "out.mkv"); err != nil {
		log.Fatal(err)
	}
}
