package main

import (
	"log"

	"github.com/peterstace/grayt/tracer"
)

func main() {

	mov := tracer.Movie{
		FrameCount: 512,
		CameraFactory: func(t float64) tracer.Camera {
			return tracer.NewRectilinearCamera(
				tracer.CameraConfig{
					Location:      tracer.NewVect(0, 10, 5),
					ViewDirection: tracer.NewVect(0, -1, -2),
					UpDirection:   tracer.NewVect(0, 1, 0),
					FieldOfView:   1.5,
					FocalLength:   10.0,
					FocalRatio:    1000.0,
				},
			)
		},
		GeometriesFactory: func(t float64) []tracer.Geometry {
			return []tracer.Geometry{
				tracer.NewPlane(tracer.NewVect(0, 0, 0), tracer.NewVect(0, 1, 0)),
				tracer.NewPlane(tracer.NewVect(0, -0.5, -10), tracer.NewVect(0.2, 1, 0)),
				tracer.NewPlane(tracer.NewVect(0, -0.5, -10), tracer.NewVect(-0.2, 1, 0)),
				tracer.NewPlane(tracer.NewVect(0, -0.5, -10), tracer.NewVect(0, 1, 0.2)),
				tracer.NewPlane(tracer.NewVect(0, -0.5, -10), tracer.NewVect(0, 1, -0.2)),
				tracer.NewSphere(tracer.NewVect(0, 1, -10), t*2),
			}
		},
		LightsFactory: func(t float64) []tracer.Light {
			return []tracer.Light{
				tracer.Light{Location: tracer.NewVect(0, 10, -10), Radius: 0.5, Intensity: 100},
			}
		},
	}

	if err := tracer.TraceMovie(mov, "out"); err != nil {
		log.Fatal(err)
	}
}
