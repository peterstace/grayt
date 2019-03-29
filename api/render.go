package api

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/peterstace/grayt/colour"
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

	url := fmt.Sprintf(
		"http://%s/trace?scene_name=%s&px_wide=%d&px_high=%d&depth=%d",
		r.workerAddr, r.scene, r.pxWide, r.pxHigh, depth,
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error requesting  work: %v", err)
		return
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			body = []byte("couldn't ready body")
		}
		log.Printf("worker result with non-200 status: %s", string(body))
		return
	}

	// TODO: should be able to reuse the pixel grid to save on allocations
	pixels := r.pxWide * r.pxHigh
	unitOfWork := pixelGrid{
		r.pxWide, r.pxHigh,
		make([]colour.Colour, pixels),
	}
	if err := binary.Read(resp.Body, binary.BigEndian, &unitOfWork.pixels); err != nil {
		log.Printf("could not read from worker response body: %v", err)
		return
	}
	if n, err := resp.Body.Read([]byte{0}); n != 0 || err != io.EOF {
		log.Printf("more bytes in response body than expected")
		return
	}

	r.acc.merge(&unitOfWork, depth)
	r.monitor.addPoint(time.Now(), pixels*depth)
}
