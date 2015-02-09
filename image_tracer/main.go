package main

import (
	"encoding/json"
	"flag"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/peterstace/grayt/tracer"
)

func main() {

	iFlag := flag.String("i", "", "input JSON file")
	oFlag := flag.String("o", "", "output JPG file")
	flag.Parse()
	if *iFlag == "" || *oFlag == "" {
		flag.PrintDefaults()
		return
	}

	var scene tracer.Scene
	if buf, err := ioutil.ReadFile(*iFlag); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(buf, &scene); err != nil {
		log.Fatal(err)
	}

	img := tracer.TraceImage([]tracer.Scene{scene})

	if out, err := os.Create(*oFlag); err != nil {
		log.Fatal(err)
	} else if err := jpeg.Encode(out, img, nil); err != nil {
		log.Fatal(err)
	}
}
