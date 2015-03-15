package main

import (
	"log"

	"github.com/peterstace/grayt/movie"
	"github.com/peterstace/grayt/tracer"
	"github.com/peterstace/grayt/vect"
)

func main() {

	mov := movie.Movie{
		Frames: 256,
		Camera: movie.ConstCamera(tracer.NewRectilinearCamera(
			tracer.CameraConfig{
				Location:      vect.New(0, 10, 5),
				ViewDirection: vect.New(0, -1, -2),
				UpDirection:   vect.New(0, 1, 0),
				FieldOfView:   1.5,
				FocalLength:   10.0,
				FocalRatio:    1000.0,
			})),
		Geometries: []movie.Geometry{

			movie.ConstGeometry(tracer.NewPlane(vect.New(0, 0, 0), vect.New(0, 1, 0))),

			movie.ConstGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(0.2, 1, 0))),
			movie.ConstGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(-0.2, 1, 0))),
			movie.ConstGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(0, 1, 0.2))),
			movie.ConstGeometry(tracer.NewPlane(vect.New(0, -0.5, -10), vect.New(0, 1, -0.2))),

			movie.ConstGeometry(tracer.NewSphere(vect.New(0, 1, -10), 1)),
		},

		Lights: []movie.Light{
			movie.ConstLight(tracer.Light{Location: vect.New(0, 10, -10), Radius: 0.5, Intensity: 100}),
		},
	}

	if err := movie.TraceMovie(mov, "out.mkv"); err != nil {
		log.Fatal(err)
	}
}
