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

	var out string

	flag.StringVar(&out, "o", "", "output file (must end in .png)")
	flag.Parse()

	if !strings.HasSuffix(out, ".png") {
		log.Fatalf(`%q does not end in ".png"`, out)
	}

	scene := CornellBoxStandard()

	acc := grayt.NewAccumulator(300, 300)
	grayt.PathTracer(scene, acc, 1000)
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
