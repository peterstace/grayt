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
		out            string
		spp            int
		cv             float64
		pxWide, pxHigh int
	)

	// Set up flags.
	flag.StringVar(&out, "o", "",
		"output file (must end in .png)")
	flag.IntVar(&spp, "spp", 0,
		"samples per pixel stopping point")
	flag.Float64Var(&cv, "cv", 0.0,
		"neighbourhood CV (coefficient of variation) stopping point")
	flag.IntVar(&pxWide, "w", 0, "width in pixels")
	flag.IntVar(&pxHigh, "h", 0, "height in pixels")
	flag.Parse()

	// Validate and interpret flags.
	if !strings.HasSuffix(out, ".png") {
		flag.Usage()
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
	if pxWide == 0 || pxHigh == 0 {
		flag.Usage()
		log.Fatal("width and height must be set")
	}

	// Load scene. TODO: load from file.
	scene := CornellBox()

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

	smoothedSamplesPerSecond := expSmoothedVar{alpha: 0.1}

	var world grayt.World
	world.AddEntities(scene.Entities)

	for !mode.stop() {

		startIterationTime := time.Now()

		grayt.TracerImage(scene.Camera, world, acc)
		iteration++
		cv := acc.NeighbourCoefficientOfVariation()
		mode.finishSample(cv)

		samplesPerSecondInIteration := float64(totalPx) / time.Now().Sub(startIterationTime).Seconds()
		smoothedSamplesPerSecond.next(samplesPerSecondInIteration)

		totalSamples := mode.estSamplesPerPixelRequired()
		totalSamplesStr := "??"
		eta := "??"
		if totalSamples >= 0 {
			totalSamplesStr = fmt.Sprintf("%d", totalSamples)
			etaSeconds := float64((totalSamples-iteration)*totalPx) / smoothedSamplesPerSecond.value
			eta = fmt.Sprintf("%v", time.Duration(etaSeconds)*time.Second)
		}

		log.Printf("Sample=%d/%s, Samples/sec=%.2e CV=%.4f ETA=%s\n",
			iteration, totalSamplesStr, samplesPerSecondInIteration, cv, eta)
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
	completed        int
	cvDeltaPerSample expSmoothedVar
}

func (u *untilRelativeStdDevBelowThreshold) estSamplesPerPixelRequired() int {
	if u.cvDeltaPerSample.value < 0 {
		return -1
	}
	more := (u.currentCV - u.threshold) / u.cvDeltaPerSample.value
	return u.completed + int(more)
}

func (u *untilRelativeStdDevBelowThreshold) finishSample(relStdDev float64) {
	u.currentCV, u.previousCV = relStdDev, u.currentCV
	u.completed++
	u.cvDeltaPerSample.alpha = 0.3
	u.cvDeltaPerSample.next(u.previousCV - u.currentCV)
}

func (u *untilRelativeStdDevBelowThreshold) stop() bool {
	return u.cvDeltaPerSample.value > 0 && u.currentCV < u.threshold
}

type expSmoothedVar struct {
	alpha float64
	value float64
}

func (v *expSmoothedVar) next(n float64) {
	if v.value == 0 {
		v.value = n
	} else {
		v.value = v.alpha*n + (1.0-v.alpha)*v.value
	}
}
