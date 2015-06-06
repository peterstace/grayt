package main

import (
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/peterstace/grayt"
)
import (
	"flag"
	"fmt"
)

func main() {

	var (
		out  string
		spp  int
		nrsd float64
	)

	flag.StringVar(&out, "o", "", "output file (must end in .png)")
	flag.IntVar(&spp, "s", 0, "samples per pixel stopping point")
	flag.Float64Var(&nrsd, "d", 0.0, "neighbourhood relative std dev stopping point")
	flag.Parse()

	if !strings.HasSuffix(out, ".png") {
		log.Fatalf(`%q does not end in ".png"`, out)
	}
	if (spp == 0 && nrsd == 0) || (spp != 0 && nrsd != 0) {
		log.Fatalf(`exactly 1 of s and d must be set`)
	}

	scene := CornellBox()

	const (
		pxWide  = 300
		pxHigh  = 300
		totalPx = pxWide * pxHigh
	)
	acc := grayt.NewAccumulator(pxWide, pxHigh)

	startTime := time.Now()
	iteration := 0
	var prevRSD, currentRSD float64
	var downRSDCount int
	for {
		iteration++
		samplesPerSecond := float64(iteration*totalPx) / time.Now().Sub(startTime).Seconds()
		timeRemaining := "??"
		totalSamples := "??"
		if spp != 0 {
			timeRemaining = fmt.Sprintf("%v",
				time.Duration(((spp-iteration)*totalPx)/int(samplesPerSecond))*time.Second)
			totalSamples = fmt.Sprintf("%d", spp)
		}
		grayt.TracerImage(scene, acc)
		prevRSD, currentRSD = currentRSD, acc.NeighbourRelativeStdDev()
		log.Printf("Sample=%d/%s, Samples/sec=%.2e NRSD=%.4f ETA=%s\n",
			iteration, totalSamples, samplesPerSecond, currentRSD, timeRemaining)
		if currentRSD < prevRSD {
			downRSDCount++
		} else {
			downRSDCount = 0
		}
		if downRSDCount > 5 && currentRSD < nrsd {
			break
		}
		if spp != 0 && iteration == spp {
			break
		}
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
