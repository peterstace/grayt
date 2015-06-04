package main

import (
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/peterstace/grayt"
)
import "flag"

func main() {

	var (
		out string
		spp int
	)

	flag.StringVar(&out, "o", "", "output file (must end in .png)")
	flag.IntVar(&spp, "s", 10, "samples per pixel")
	flag.Parse()

	if !strings.HasSuffix(out, ".png") {
		log.Fatalf(`%q does not end in ".png"`, out)
	}

	scene := CornellBox()

	acc := grayt.NewAccumulator(300, 300)
	grayt.TracerImage(scene, acc, grayt.Quality{spp})
	img := acc.ToImage(1.0)

	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
