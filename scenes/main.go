package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/peterstace/grayt"
)

func main() {

	scenes := map[string]func() grayt.Scene{
		"cornell": CornellBox,
	}

	var s string

	flag.StringVar(&s, "s", "", "scene description (one of "+string("")+")")
	flag.Parse()

	fn, ok := scenes[s]
	if !ok {
		flag.Usage()
		log.Printf("scene %q not found", s)
		return
	}

	buf, err := json.Marshal(fn())
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Print(string(buf))
	fmt.Print("\n")
}
