package main

import (
	"log"
	"time"

	"github.com/peterstace/grayt/movie"
	"github.com/peterstace/grayt/tracer"
	"github.com/peterstace/grayt/vect"
)

func main() {

	mov := movie.Movie{
		Duration: time.Second * 17,
		Camera: movie.StillCamera(tracer.NewRectilinearCamera(
			tracer.CameraConfig{
				Location:      vect.New(0, 10, 5),
				ViewDirection: vect.New(0, -1, -2),
				UpDirection:   vect.New(0, 1, 0),
				FieldOfView:   1.5,
				FocalLength:   10.0,
				FocalRatio:    1000.0,
			})),
		Geometries: []movie.GeometryFn{

			movie.StillGeometry(tracer.NewPlane(vect.New(0, 0, 0), vect.New(0, 1, 0))),

			movie.StillGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(0.2, 1, 0))),
			movie.StillGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(-0.2, 1, 0))),
			movie.StillGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(0, 1, 0.2))),
			movie.StillGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(0, 1, -0.2))),

			movie.StillGeometry(tracer.NewSphere(vect.New(0, 1, -10), 1)),
		},

		Lights: []movie.LightFn{
			movie.StillLight(tracer.Light{Location: vect.New(0, 10, -10), Radius: 0.5, Intensity: 100}),
		},
	}

	if err := movie.TraceMovie(mov, "out.mkv"); err != nil {
		log.Fatal(err)
	}
}
