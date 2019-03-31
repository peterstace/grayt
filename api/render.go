package api

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/scene/library"
	"github.com/peterstace/grayt/trace"
)

type render struct {
	scene   string
	pxWide  int
	pxHigh  int
	created time.Time

	cnd            *sync.Cond
	desiredWorkers int
	actualWorkers  int

	acc *accumulator

	monitor rateMonitor

	workerAddr string
}

func (r *render) orchestrateWork() {
	for {
		r.cnd.L.Lock()
		for r.actualWorkers >= r.desiredWorkers {
			r.cnd.Wait()
		}
		r.actualWorkers++
		go r.work()
		r.cnd.L.Unlock()
	}
}

func (r *render) work() {
	defer func() {
		r.cnd.L.Lock()
		r.actualWorkers--
		r.cnd.L.Unlock()
		r.cnd.Signal()
	}()

	const depth = 30

	sceneFn, ok := library.Lookup(r.scene)
	if !ok {
		panic(fmt.Sprintf("unknown: %v", r.scene))
	}
	scn := trace.BuildScene(sceneFn())
	accel := trace.NewGrid(4, scn.Objects)

	var buf bytes.Buffer
	traceLayer(&buf, r.pxWide, r.pxHigh, depth, accel, scn.Camera)

	// TODO: should be able to reuse the pixel grid to save on allocations
	pixels := r.pxWide * r.pxHigh
	unitOfWork := pixelGrid{
		r.pxWide, r.pxHigh,
		make([]colour.Colour, pixels),
	}
	if err := binary.Read(&buf, binary.BigEndian, &unitOfWork.pixels); err != nil {
		log.Printf("could not read from worker response body: %v", err)
		return
	}
	if n, err := buf.Read([]byte{0}); n != 0 || err != io.EOF {
		log.Printf("more bytes in response body than expected")
		return
	}

	r.acc.merge(&unitOfWork, depth)
	r.monitor.addPoint(time.Now(), pixels*depth)
}
