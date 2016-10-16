package grayt

import (
	"flag"
	"fmt"
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
		case <-time.After(100 * time.Millisecond):
			cli.update(int(atomic.LoadUint64(completed)))
		case img := <-img:
			cli.finished()
			output := fmt.Sprintf("%s[%s]_%dx%d_q%d.png",
				baseName, scene.hash(), *pxWide, *pxHigh, *quality)
			outFile, err := os.Create(output)
			if err != nil {
				log.Fatal(err)
			}
			err = png.Encode(outFile, img)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}
	}
}
