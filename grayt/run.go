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
	"reflect"
	"sync/atomic"
	"time"
)

var (
	pxWide  = flag.Int("w", 640, "width in pixels")
	quality = flag.Int("q", 10, "quality (samples per pixel)")
	verbose = flag.Bool("v", false, "verbose model")
	output  = flag.String("o", "", "output file override")
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

	pxHigh := *pxWide * scene.Camera.aspectHigh / scene.Camera.aspectWide

	accel := newGrid(2, scene.Objects)
	log.Printf("Origin:     %#v", accel.minBound)
	log.Printf("Stride:     %#v", accel.stride)
	log.Printf("Resolution: %#v", accel.resolution)
	for i := range accel.data {
		log.Printf("Data:       %v", accel.data[i])
	}
	//accel := newListAccelerationStructure(scene.Objects)

	cam := newCamera(scene.Camera)
	img := make(chan image.Image)
	completed := new(uint64)
	go func() {
		img <- traceImage(*pxWide, pxHigh, accel, cam, *quality, completed)
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
				*output = fmt.Sprintf("%s[%s]_%dx%d_q%d.png",
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

	// Calculate hash.
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	binary.Write(h, binary.LittleEndian, s.Camera)
	for _, o := range s.Objects {
		h.Write([]byte(reflect.TypeOf(o).String()))
		binary.Write(h, binary.LittleEndian, o)
	}

	// Calculate base64 encoded hash.
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
