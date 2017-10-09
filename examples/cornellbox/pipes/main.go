package main

import (
	"math/rand"

	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

func scene() Scene {
	pipes := Group()
	for i := 0; i < 100; i++ {
		a := rand.Float64()*0.8 + 0.1
		b := rand.Float64()*0.8 + 0.1
		c := rand.Float64()*0.8 + 0.1
		d := rand.Float64()*0.8 + 0.1
		pipes = Group(pipes, Pipe(
			Vect(0, a, -b),
			Vect(1, c, -d),
			0.01,
		))
	}
	return Scene{
		Camera: Cam(1.3),
		Objects: Group(
			Floor,
			Ceiling,
			BackWall,
			LeftWall.With(ColourRGB(Red)),
			RightWall.With(ColourRGB(Green)),
			CeilingLight().With(Emittance(5.0)),
			pipes,
		),
	}
}

func main() {
	Run("cornellbox_classic", scene())
}
