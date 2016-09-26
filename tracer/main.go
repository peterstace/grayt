package main

import (
	"flag"
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

	f, err := os.Open(*inputFlag)
	if err != nil {
		log.Fatal(err)
	}

	s, err := scene.ReadFrom(f)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(s)
}
