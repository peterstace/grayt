package grayt

import (
	"image/png"
	"log"
	"os"
)

// Runner is a convenience struct to help run grayt from a main() function.
type Runner struct {
	PxWide, PxHigh int
	BaseName       string
	Quality        float64
}

func NewRunner() *Runner {
	return &Runner{
		PxWide:   640,
		PxHigh:   480,
		BaseName: "default",
		Quality:  10,
	}
}

func (r *Runner) Run(scene Scene) {

	var world World
	world.AddEntities(scene.Entities)

	camera, err := NewCamera(scene.CameraConfig)
	r.checkErr(err)

	acc := NewAccumulator(r.PxWide, r.PxHigh)
	for i := 0; i < int(r.Quality); i++ {
		log.Print(i)
		TracerImage(camera, world, acc)
	}
	img := acc.ToImage(1.0) // XXX should be configurable
	f, err := os.Create(r.BaseName + ".png")
	r.checkErr(err)
	defer f.Close()
	err = png.Encode(f, img)
	r.checkErr(err)
}

func (r *Runner) checkErr(err error) {
	if err != nil {
		log.Fatal("Fatal: ", err)
	}
}
