package grayt

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

// Run runs the grayt framework. The scene can be loaded using the provided
// function. It should be called from your main function.
func Run(fn func(*API)) {
	api := &API{
		aspectRatio: [2]int{4, 3},
	}
	fn(api)

	var fl flags
	fl.parse()

	img := make(chan image.Image)
	completed := new(int64)
	go func() {
		img <- traceImage(fl, api, completed)
	}()

	pxHigh := fl.pxWide * api.aspectRatio[1] / api.aspectRatio[0]
	total := fl.pxWide * pxHigh * fl.quality
	cli := newCLI(total)

	for {
		select {
		case <-time.After(updateInterval):
			cli.update(int(atomic.LoadInt64(completed)))
		case img := <-img:
			cli.finished()
			output := fmt.Sprintf("%s_%dx%d_q%d.png",
				time.Now().Format("20060102-150405"),
				fl.pxWide, pxHigh, fl.quality,
			)
			outFile, err := os.Create(output)
			if err != nil {
				log.Fatal(err)
			}
			enc := png.Encoder{CompressionLevel: png.BestCompression}
			err = enc.Encode(outFile, img)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
}

type flags struct {
	pxWide  int
	quality int
	workers int
}

func (f *flags) parse() {
	flag.IntVar(&f.pxWide, "w", 640, "width in pixels")
	flag.IntVar(&f.quality, "q", 10, "quality (samples per pixel)")
	flag.IntVar(&f.pxWide, "j", runtime.GOMAXPROCS(0), "number of worker goroutines")
	flag.Parse()
}
