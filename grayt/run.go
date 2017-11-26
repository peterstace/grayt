package grayt

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"image"
	"image/png"
	"log"
	"os"
	"sync/atomic"
	"time"
)

var scenes = map[string]func() Scene{}

func Register(name string, fn func() Scene) {
	scenes[name] = fn
}

func RunScene() error {
	if fn, ok := scenes[Config.Scene]; ok {
		Run(Config.Scene, fn())
		return nil
	} else {
		return fmt.Errorf("could not find scene %q", Config.Scene)
	}
}

// Run should be the single call made from main().
func Run(baseName string, scene Scene) {
	if Config.Verbose {
		fmt.Printf("Camera: %v\n", scene.Camera)
		for i, o := range scene.Objects {
			fmt.Printf("Object %d:\n%v\n", i, o)
		}
	}

	pxHigh := scene.Camera.pxHigh(Config.PxWide)
	accum := initAccumulator(Config.PxWide, pxHigh)
	img := make(chan image.Image)
	completed := uint64(accum.count * (Config.PxWide) * pxHigh)
	go func() {
		img <- TraceImage(Config.PxWide, scene, Config.Quality, Config.NumWorkers, accum, &completed)
	}()

	total := Config.PxWide * pxHigh * Config.Quality
	cli := newCLI(total)

	for {
		select {
		case <-time.After(updateInterval):
			cli.update(int(atomic.LoadUint64(&completed)))
		case img := <-img:
			cli.finished()
			if Config.Output == "" {
				Config.Output = fmt.Sprintf("%s_%s_%s_%dx%d_q%d.png",
					time.Now().Format("20060102-150405"),
					baseName, hashScene(scene), Config.PxWide, pxHigh, Config.Quality)
			}
			outFile, err := os.Create(Config.Output)
			if err != nil {
				log.Fatal(err)
			}
			enc := png.Encoder{CompressionLevel: png.BestCompression}
			err = enc.Encode(outFile, img)
			if err != nil {
				log.Fatal(err)
			}
			os.Remove("checkpoint") // Ignore error.
			os.Exit(0)
		}
	}
}

func hashScene(s Scene) string {
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	fmt.Fprintf(h, "%+v", s)

	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.RawURLEncoding, &buf)
	binary.Write(
		enc,
		binary.LittleEndian,
		h.Sum64(),
	)
	enc.Close()
	return buf.String()
}

func initAccumulator(wide, high int) *accumulator {
	accum := new(accumulator)
	n := wide * high
	if ok, err := accum.load(); err != nil {
		log.Fatal("Could not load checkpoint:", err)
	} else if ok {
		// TODO: Check other properties, such as appropriate flags and scene hash.
		if accum.wide != wide || accum.high != high {
			log.Fatalf("Checkpoint size doesn't match settings: checkpoint=%dx%d settings=%dx%d",
				accum.wide, accum.high, wide, high)
		}
	} else {
		accum.pixels = make([]Colour, n)
		accum.wide = wide
		accum.high = high
	}
	return accum
}
