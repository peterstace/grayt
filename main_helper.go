package grayt

import (
	"fmt"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Runner is a convenience struct to help run grayt from a main() function.
type Runner struct {
	PxWide, PxHigh int
	BaseName       string
	Quality        int
}

func NewRunner() *Runner {
	return &Runner{
		PxWide:   640,
		PxHigh:   480,
		BaseName: "default",
		Quality:  10,
	}
}

func (r *Runner) Run(s Scene) {
	fmt.Println(r.outputFilename())
	strat := strategy{}
	img := strat.traceImage(r.PxWide, r.PxHigh, s, r.Quality)
	f, err := os.Create(r.outputFilename())
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

func (r *Runner) outputFilename() string {
	return fmt.Sprintf("%s_%dx%d_Q%d.png",
		r.BaseName, r.PxWide, r.PxHigh, r.Quality)
}
