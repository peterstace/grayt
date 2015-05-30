package main

import (
	"image/png"
	"os"
	"strings"

	"github.com/peterstace/grayt"
)
import "flag"

func main() {

	var out string

	flag.StringVar(&out, "o", "", "output file (must end in .jpeg, .jpg, or .png)")
	flag.Parse()

	if getOutType(out) == outTypeUnknown {
		flag.Usage()
		return
	}

	scene := CornellBoxStandard()

	acc := grayt.NewAccumulator(300, 300)
	grayt.RayTracer(scene, acc)
	img := acc.ToImage(1.0)

	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	switch getOutType(out) {
	case outTypePNG:
		if err := png.Encode(f, img); err != nil {
			panic(err)
		}
	case outTypeJPEG:
		panic("unimplemented") // XXX
	case outTypeUnknown:
		panic("unknown outType should have been rejected during flag validation")
	}
}

type outType int

const (
	outTypeUnknown outType = iota
	outTypePNG
	outTypeJPEG
)

func getOutType(out string) outType {
	switch {
	case strings.HasSuffix(out, ".jpeg") || strings.HasSuffix(out, ".jpg"):
		return outTypeJPEG
	case strings.HasSuffix(out, ".png"):
		return outTypePNG
	default:
		return outTypeUnknown
	}
}
