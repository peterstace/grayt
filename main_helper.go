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

	world := newWorld(scene.Entities)

	acc := newAccumulator(r.PxWide, r.PxHigh)
	for i := 0; i < int(r.Quality); i++ {
		log.Print(i)
		TracerImage(scene.Camera, world, acc)
	}
	img := acc.toImage(1.0) // XXX should be configurable
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
