package main

import (
	"flag"
	"image/png"
	"log"
	"os"

	"github.com/peterstace/grayt/scene"
)

func main() {

	inputFlag := flag.String("f", "", "input file")
	outputFlag := flag.String("o", "", "output file")
	flag.Parse()

	if *inputFlag == "" {
		log.Fatal("Input file not specified")
	}
	if *outputFlag == "" {
		log.Fatal("Output file not specified.")
	}

	inFile, err := os.Open(*inputFlag)
	if err != nil {
		log.Fatal(err)
	}

	s, err := scene.ReadFrom(inFile)
	if err != nil {
		log.Fatal(err)
	}

	img := traceImage(s)
	outFile, err := os.Create(*outputFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(outFile, img)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(s)
}
