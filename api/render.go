package api

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
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
}

func (r *render) orchestrateWork() {
	// TODO: Could just get lock above loop and then never release?
	// TODO: Add some sort of backoff.
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
		log.Printf("too many requests, trying again")
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

	// TODO: read body and put into accumulator
}

func (r *render) image() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, r.pxWide, r.pxHigh))
	for x := 0; x < r.pxWide; x++ {
		for y := 0; y < r.pxHigh; y++ {
			img.Set(x, y, color.Gray{0xff})
		}
	}
	return img
}
