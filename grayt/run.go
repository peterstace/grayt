package grayt

import (
	"image"
	"image/png"
	"log"
	"os"
	"sync/atomic"
	"time"
)

const pxWide = 128
const pxHigh = 128
const quality = int(1e5)

// Run should be the single call made from main().
func Run(baseName string, scene Scene) {

	tris := convertTriangles(scene.Triangles)
	accel := newAccelerationStructure(tris)
	cam := newCamera(scene.Camera)
	img := make(chan image.Image)
	completed := new(uint64)
	go func() {
		img <- traceImage(pxWide, pxHigh, accel, cam, quality, completed)
	}()

	total := pxWide * pxHigh * quality
	cli := newCLI(total)

	for {
		select {
		case <-time.After(100 * time.Millisecond):
			cli.update(int(atomic.LoadUint64(completed)))
		case img := <-img:
			cli.finished()
			output := baseName + ".png"
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
