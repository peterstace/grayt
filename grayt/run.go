package grayt

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"hash/crc64"
	"image"
	"image/png"
	"log"
	"os"
	"sync/atomic"
	"time"
)

var (
	pxWide     = flag.Int("w", 640, "width in pixels")
	quality    = flag.Int("q", 10, "quality (samples per pixel)")
	verbose    = flag.Bool("v", false, "verbose model")
	output     = flag.String("o", "", "output file override")
	numWorkers = flag.Int("j", 1, "number of worker goroutines")
)

// Run should be the single call made from main().
func Run(baseName string, scene Scene) {

	flag.Parse()

	if *verbose {
		fmt.Printf("Camera: %v\n", scene.Camera)
		for i, o := range scene.Objects {
			fmt.Printf("Object %d:\n%v\n", i, o)
		}
	}

	pxHigh := scene.Camera.pxHigh(*pxWide)

	img := make(chan image.Image)
	completed := new(uint64)
	go func() {
		img <- TraceImage(*pxWide, scene, *quality, *numWorkers, completed)
	}()

	total := *pxWide * pxHigh * *quality
	cli := newCLI(total)

	for {
		select {
		case <-time.After(updateInterval):
			cli.update(int(atomic.LoadUint64(completed)))
		case img := <-img:
			cli.finished()
			if *output == "" {
				*output = fmt.Sprintf("%s_%s_%s_%dx%d_q%d.png",
					time.Now().Format("20060102-150405"),
					baseName, hashScene(scene), *pxWide, pxHigh, *quality)
			}
			outFile, err := os.Create(*output)
			if err != nil {
				log.Fatal(err)
			}
			enc := png.Encoder{CompressionLevel: png.BestCompression}
			err = enc.Encode(outFile, img)
			if err != nil {
				log.Fatal(err)
			}
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
