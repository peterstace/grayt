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

// Run should be the single call made from main().
func Run(baseName string, scene Scene) {

	pxWide := flag.Int("w", 640, "width in pixels")
	pxHigh := flag.Int("h", 480, "height in pixels")
	quality := flag.Int("q", 10, "quality (samples per pixel)")
	flag.Parse()

	tris := convertTriangles(scene.Triangles)
	accel := newAccelerationStructure(tris)
	cam := newCamera(scene.Camera)
	img := make(chan image.Image)
	completed := new(uint64)
	go func() {
		img <- traceImage(*pxWide, *pxHigh, accel, cam, *quality, completed)
	}()

	total := *pxWide * *pxHigh * *quality
	cli := newCLI(total)

	for {
		select {
		case <-time.After(updateInterval):
			cli.update(int(atomic.LoadUint64(completed)))
		case img := <-img:
			cli.finished()
			output := fmt.Sprintf("%s[%s]_%dx%d_q%d.png",
				baseName, hashScene(scene), *pxWide, *pxHigh, *quality)
			outFile, err := os.Create(output)
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
	h.Write([]byte(s.Camera.String()))
	for _, t := range s.Triangles {
		h.Write([]byte(t.String()))
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
