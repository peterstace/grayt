package api

import (
	"encoding/binary"
	"fmt"
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

	backoffMu sync.Mutex
	backoff   time.Duration
}

func (r *render) orchestrateWork() {
	// TODO: Could just get lock above loop and then never release?
	// TODO: Add some sort of backoff.
	for {
		r.cnd.L.Lock()
		for r.actualWorkers >= r.desiredWorkers {
			r.cnd.Wait()
		}
		r.sleepForBackoff()
		r.actualWorkers++
		go r.work()
		r.cnd.L.Unlock()
	}
}

func (r *render) sleepForBackoff() {
	r.backoffMu.Lock()
	b := r.backoff
	r.backoffMu.Unlock()
	time.Sleep(b)
}

func (r *render) work() {
	defer func() {
		r.cnd.L.Lock()
		r.actualWorkers--
		r.cnd.L.Unlock()
		r.cnd.Signal()
	}()

	// TODO: allow URL base to be configurable
	url := fmt.Sprintf(
		"http://worker:80/trace?scene_name=%s&px_wide=%d&px_high=%d",
		r.scene, r.pxWide, r.pxHigh,
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error requesting  work: %v", err)
		return
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		r.backoffMu.Lock()
		r.backoff = 2 * (time.Millisecond + r.backoff)
		backoff := r.backoff
		r.backoffMu.Unlock()
		log.Printf("too many requests, backoff: %s", backoff)
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

	r.backoffMu.Lock()
	r.backoff = 0
	r.backoffMu.Unlock()

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

	r.acc.merge(&unitOfWork)
}
