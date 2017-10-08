package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/peterstace/grayt/engine"
	_ "github.com/peterstace/grayt/scenes"
)

type options struct {
	scene      string
	listScenes bool
	pxWide     int
	quality    int
	workers    int
}

func (o *options) parse() {
	flag.StringVar(&o.scene, "s", "", "scene to render")
	flag.BoolVar(&o.listScenes, "l", false, "list available scenes")
	flag.IntVar(&o.pxWide, "w", 640, "width in pixels")
	flag.IntVar(&o.quality, "q", 10, "quality (samples per pixel)")
	flag.IntVar(&o.workers, "j", runtime.GOMAXPROCS(0), "number of worker goroutines")
	flag.Parse()
}

func main() {
	var opt options
	opt.parse()
	if opt.listScenes {
		scenes := engine.SceneList()
		fmt.Printf("%d scenes available:\n", len(scenes))
		for _, s := range scenes {
			fmt.Println(s)
		}
		return
	}
	sceneFn, ok := engine.LookupScene(opt.scene)
	if !ok {
		fmt.Fprintf(os.Stderr, "Could not find registered scene %q\n", opt.scene)
		os.Exit(1)
	}
	if err := run(sceneFn, opt); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(scene func(*engine.API), opt options) error {
	img := make(chan image.Image)
	var status engine.Status
	go func() {
		img <- engine.TraceImage(opt.pxWide, scene, opt.quality, &status)
	}()
	var display cli
	for {
		select {
		case <-time.After(updateInterval):
			display.update(
				int(atomic.LoadInt64(&status.Done)),
				int(atomic.LoadInt64(&status.Total)),
			)
		case img := <-img:
			display.finished(
				int(atomic.LoadInt64(&status.Total)),
			)
			imgSz := img.Bounds().Size()
			output := fmt.Sprintf("%s_%dx%d_q%d.png",
				time.Now().Format("20060102-150405"),
				imgSz.X, imgSz.Y, opt.quality,
			)
			outFile, err := os.Create(output)
			if err != nil {
				return err
			}
			enc := png.Encoder{CompressionLevel: png.BestCompression}
			err = enc.Encode(outFile, img)
			if err != nil {
				return err
			}
			return nil
		}
	}
}
