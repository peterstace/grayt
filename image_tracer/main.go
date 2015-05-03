package main

import (
	"log"

	"github.com/peterstace/grayt"
)

func main() {

	mov := grayt.Movie{
		FrameCount: 128,
		CameraFactory: func(t float64) grayt.Camera {
			return grayt.NewRectilinearCamera(
				grayt.CameraConfig{
					Location:      grayt.Vect{0, 10, 5},
					ViewDirection: grayt.Vect{0, -1, -2},
					UpDirection:   grayt.Vect{0, 1, 0},
					FieldOfView:   1.5,
					FocalLength:   10.0,
					FocalRatio:    1000.0,
				},
			)
		},
		GeometriesFactory: func(t float64) []grayt.Geometry {
			return []grayt.Geometry{
				grayt.NewPlane(grayt.Vect{0, 0, 0}, grayt.Vect{0, 1, 0}),
				grayt.NewPlane(grayt.Vect{0, -0.5, -10}, grayt.Vect{0.2, 1, 0}),
				grayt.NewPlane(grayt.Vect{0, -0.5, -10}, grayt.Vect{-0.2, 1, 0}),
				grayt.NewPlane(grayt.Vect{0, -0.5, -10}, grayt.Vect{0, 1, 0.2}),
				grayt.NewPlane(grayt.Vect{0, -0.5, -10}, grayt.Vect{0, 1, -0.2}),
				grayt.NewSphere(grayt.Vect{0, 1, -10}, t*2),
			}
		},
		LightsFactory: func(t float64) []grayt.Light {
			return []grayt.Light{
				grayt.Light{Location: grayt.Vect{0, 10, -10}, Radius: 0.5, Intensity: 100},
			}
		},
	}

	if err := grayt.TraceMovie(mov, "out.mkv"); err != nil {
		log.Fatal(err)
	}
}
