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
		out string
		spp int
		cv  float64
	)

	// Set up flags.
	flag.StringVar(&out, "out", "",
		"output file (must end in .png)")
	flag.IntVar(&spp, "spp-limit", 0,
		"samples per pixel stopping point")
	flag.Float64Var(&cv, "cv-limit", 0.0,
		"neighbourhood CV (coefficient of variation) stopping point")
	flag.Parse()

	// Validate and interpret flags.
	if !strings.HasSuffix(out, ".png") {
		log.Fatalf(`%q does not end in ".png"`, out)
	}
	if (spp == 0 && cv == 0) || (spp != 0 && cv != 0) {
		log.Fatalf(`exactly 1 of s and d must be set`)
	}
	var mode mode
	if spp != 0 {
		mode = &fixedSamplesPerPixel{required: spp}
	} else {
		mode = &untilRelativeStdDevBelowThreshold{threshold: cv}
	}

	// Load scene. TODO: load from file.
	scene := CornellBox()

	// TODO: these should come from command line args.
	const (
		pxWide  = 300
		pxHigh  = 300
		totalPx = pxWide * pxHigh
	)
	acc := grayt.NewAccumulator(pxWide, pxHigh)

	run(mode, scene, acc)

	// TODO: image exposure should come from the scene description.
	img := acc.ToImage(1.0)

	// Write image out to file.
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
		cv := acc.NeighbourCoefficientOfVariation()
		mode.finishSample(cv)

		samplesPerSecond := float64(iteration*totalPx) / time.Now().Sub(startTime).Seconds()

		totalSamples := mode.estSamplesPerPixelRequired()
		totalSamplesStr := "??"
		eta := "??"
		if totalSamples >= 0 {
			totalSamplesStr = fmt.Sprintf("%d", totalSamples)
			etaSeconds := float64((totalSamples-iteration)*totalPx) / samplesPerSecond
			eta = fmt.Sprintf("%v", time.Duration(etaSeconds)*time.Second)
		}

		log.Printf("Sample=%d/%s, Samples/sec=%.2e CV=%.4f ETA=%s\n",
			iteration, totalSamplesStr, samplesPerSecond, cv, eta)
	}

	log.Printf("TotalTime=%s", time.Now().Sub(startTime))
}

type mode interface {
	// estSamplesPerPixelRequired is an estimation of the number of total
	// samples per pixel that will be required before the render is completed.
	estSamplesPerPixelRequired() int

	// finishSample signals to the mode that a sample has been finished. The
	// new coefficient of variation should be supplied.
	finishSample(cv float64)

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
	threshold        float64
	currentCV        float64
	previousCV       float64
	reducedCVCount   int
	completed        int
	cvDeltaPerSample float64
}

func (u *untilRelativeStdDevBelowThreshold) estSamplesPerPixelRequired() int {
	if u.reducedCVCount < 5 {
		return -1
	}
	more := (u.currentCV - u.threshold) / u.cvDeltaPerSample
	return u.completed + int(more)
}

func (u *untilRelativeStdDevBelowThreshold) finishSample(relStdDev float64) {
	u.currentCV, u.previousCV = relStdDev, u.currentCV
	if u.currentCV < u.previousCV {
		u.reducedCVCount++
	} else {
		u.reducedCVCount = 0
	}
	u.completed++
	u.cvDeltaPerSample = 0.9*u.cvDeltaPerSample + 0.1*(u.previousCV-u.currentCV)
}

func (u *untilRelativeStdDevBelowThreshold) stop() bool {
	return u.reducedCVCount >= 5 && u.currentCV < u.threshold
}
