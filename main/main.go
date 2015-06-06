package main

import (
	"image/png"
	"log"
	"os"
	"strings"
	"time"

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

	const (
		pxWide  = 300
		pxHigh  = 300
		totalPx = pxWide * pxHigh
	)
	acc := grayt.NewAccumulator(pxWide, pxHigh)
	startTime := time.Now()
	for i := 1; i <= spp; i++ {
		samplesPerSecond := float64(i*totalPx) / time.Now().Sub(startTime).Seconds()
		timeRemaining := time.Duration(((spp-i)*totalPx)/int(samplesPerSecond)) * time.Second
		log.Printf("Sample=%d/%d Sampes/sec=%.2e NRSD=%.2f%% ETA=%s\n",
			i, spp, samplesPerSecond, acc.NeighbourRelativeStdDev()*100, timeRemaining)
		grayt.TracerImage(scene, acc)
	}
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
