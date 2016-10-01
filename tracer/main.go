package main

import (
	"flag"
	"image/png"
	"log"
	"os"

	"github.com/peterstace/grayt/scene"
)

type flags struct {
	input  *string
	output *string
	pxWide *int
	pxHigh *int
}

func getFlags() flags {
	f := flags{
		input:  flag.String("f", "", "Input file"),
		output: flag.String("o", "", "Output file"),
		pxWide: flag.Int("w", 640, "Width (pixels)"),
		pxHigh: flag.Int("h", 480, "Height (pixels)"),
	}
	flag.Parse()
	if *f.input == "" || *f.output == "" || *f.pxWide <= 0 || *f.pxHigh <= 0 {
		flag.Usage()
	}
	return f
}

func main() {

	f := getFlags()

	inFile, err := os.Open(*f.input)
	if err != nil {
		log.Fatal(err)
	}

	s, err := scene.ReadFrom(inFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v", s.Camera)
	for _, t := range s.Triangles {
		log.Printf("%v", t)
	}

	tris := convertTriangles(s.Triangles)
	accel := newAccelerationStructure(tris)
	cam := newCamera(s.Camera)
	img := traceImage(*f.pxWide, *f.pxHigh, accel, cam)

	outFile, err := os.Create(*f.output)
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(outFile, img)
	if err != nil {
		log.Fatal(err)
	}
}
