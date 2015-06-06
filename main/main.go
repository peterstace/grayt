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
	var mode mode
	if spp != 0 {
		mode = &fixedSamplesPerPixel{required: spp}
	} else {
		mode = &untilRelativeStdDevBelowThreshold{threshold: nrsd}
	}

	scene := CornellBox()

	const (
		pxWide  = 300
		pxHigh  = 300
		totalPx = pxWide * pxHigh
	)
	acc := grayt.NewAccumulator(pxWide, pxHigh)

	run(mode, scene, acc)

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

func run(mode mode, scene grayt.Scene, acc grayt.Accumulator) {

	wide, high := acc.Dimensions()
	totalPx := wide * high

	startTime := time.Now()
	iteration := 0

	for !mode.stop() {

		grayt.TracerImage(scene, acc)
		iteration++
		rsd := acc.NeighbourRelativeStdDev()
		mode.finishSample(rsd)

		samplesPerSecond := float64(iteration*totalPx) / time.Now().Sub(startTime).Seconds()

		totalSamples := mode.estSamplesPerPixelRequired()
		totalSamplesStr := "??"
		eta := "??"
		if totalSamples >= 0 {
			totalSamplesStr = fmt.Sprintf("%d", totalSamples)
			etaSeconds := float64((totalSamples-iteration)*totalPx) / samplesPerSecond
			eta = fmt.Sprintf("%v", time.Duration(etaSeconds)*time.Second)
		}

		log.Printf("Sample=%d/%s, Samples/sec=%.2e RSD=%.4f ETA=%s\n",
			iteration, totalSamplesStr, samplesPerSecond, rsd, eta)
	}
}

type mode interface {
	// estSamplesPerPixelRequired is an estimation of the number of total
	// samples per pixel that will be required before the render is completed.
	estSamplesPerPixelRequired() int

	// finishSample signals to the mode that a sample has been finished. The
	// new relative std dev should be supplied.
	finishSample(relativeStdDev float64)

	// stop indicates if the render is complete.
	stop() bool
}

type fixedSamplesPerPixel struct {
	required  int
	completed int
}

func (f *fixedSamplesPerPixel) estSamplesPerPixelRequired() int {
	return f.required
}

func (f *fixedSamplesPerPixel) finishSample(float64) {
	f.completed++
}

func (f *fixedSamplesPerPixel) stop() bool {
	return f.required == f.completed
}

type untilRelativeStdDevBelowThreshold struct {
	threshold         float64
	currentRSD        float64
	previousRSD       float64
	reducedRSDCount   int
	completed         int
	rsdDeltaPerSample float64
}

func (u *untilRelativeStdDevBelowThreshold) estSamplesPerPixelRequired() int {
	if u.reducedRSDCount < 5 {
		return -1
	}
	more := (u.currentRSD - u.threshold) / u.rsdDeltaPerSample
	return u.completed + int(more)
}

func (u *untilRelativeStdDevBelowThreshold) finishSample(relStdDev float64) {
	u.currentRSD, u.previousRSD = relStdDev, u.currentRSD
	if u.currentRSD < u.previousRSD {
		u.reducedRSDCount++
	} else {
		u.reducedRSDCount = 0
	}
	u.completed++
	u.rsdDeltaPerSample = 0.9*u.rsdDeltaPerSample + 0.1*(u.previousRSD-u.currentRSD)
}

func (u *untilRelativeStdDevBelowThreshold) stop() bool {
	return u.reducedRSDCount >= 5 && u.currentRSD < u.threshold
}
