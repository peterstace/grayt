package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"sync/atomic"
	"time"

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
	var err error
	if *f.input == "" {
		err = errors.New("no input file specified")
	}
	if *f.output == "" {
		err = errors.New("no output file specified")
	}
	if *f.pxWide <= 0 {
		err = errors.New("px wide must be positive")
	}
	if *f.pxHigh <= 0 {
		err = errors.New("px high must be positive")
	}
	if err != nil {
		fmt.Printf("Error while parsing flags: %s.\n", err)
		flag.Usage()
		os.Exit(1)
	}
	return f
}

func main() {

	// (\) [XXX.XX%] [99.99Z samples/sec], ETA

	f := getFlags()

	inFile, err := os.Open(*f.input)
	if err != nil {
		log.Fatal(err)
	}

	s, err := scene.ReadFrom(inFile)
	if err != nil {
		log.Fatal(err)
	}

	tris := convertTriangles(s.Triangles)
	accel := newAccelerationStructure(tris)
	cam := newCamera(s.Camera)
	img := make(chan image.Image)
	completed := new(uint64)
	go func() {
		img <- traceImage(*f.pxWide, *f.pxHigh, accel, cam, completed)
	}()

	const quality = 100
	total := *f.pxWide * *f.pxHigh * quality
	cli := newCLI(total)

	for {
		select {
		case <-time.After(time.Second):
			cli.update(int(atomic.LoadUint64(completed)))
		case img := <-img:
			cli.update(int(atomic.LoadUint64(completed)))
			cli.finished()
			outFile, err := os.Create(*f.output)
			if err != nil {
				log.Fatal(err)
			}
			err = png.Encode(outFile, img)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}
	}
}
