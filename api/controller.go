package api

import (
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"image"
	"sync"
	"time"

	"github.com/peterstace/grayt/scene/library"
	"github.com/peterstace/grayt/trace"
	"github.com/peterstace/grayt/xmath"
)

func newController() *controller {
	return &controller{
		instances: make(map[string]*instance),
	}
}

type controller struct {
	mu        sync.Mutex
	instances map[string]*instance
}

type instance struct {
	*trace.Instance
	sceneName        string
	created          time.Time
	dim              xmath.Dimensions
	requestedWorkers int
}

type render struct {
	Scene            string    `json:"scene"`
	PxWide           int       `json:"px_wide"`
	PxHigh           int       `json:"px_high"`
	Passes           int       `json:"passes"`
	Completed        string    `json:"completed"`
	TraceRate        string    `json:"trace_rate"`
	ID               string    `json:"uuid"`
	RequestedWorkers int       `json:"requested_workers"`
	ActualWorkers    int       `json:"actual_workers"`
	Created          time.Time `json:"-"`
}

func (c *controller) getRenders() []render {
	c.mu.Lock()
	defer c.mu.Unlock()
	renders := []render{}
	for id, inst := range c.instances {
		stats := inst.GetStats()
		renders = append(renders, render{
			Scene:            inst.sceneName,
			PxWide:           inst.dim.Wide,
			PxHigh:           inst.dim.High,
			Passes:           stats.Passes,
			Completed:        displayFloat64(float64(stats.Completed)),
			TraceRate:        displayFloat64(float64(stats.TraceRateHz)) + " Hz",
			ID:               id,
			RequestedWorkers: inst.requestedWorkers,
			ActualWorkers:    stats.Workers,
			Created:          inst.created,
		})
	}
	return renders
}

func (c *controller) newRender(sceneName string, dim xmath.Dimensions) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	sceneFn, ok := library.Lookup(sceneName)
	if !ok {
		return "", fmt.Errorf("unknown scene name: %v", sceneName)
	}

	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(time.Now().Unix()))
	sum := crc64.Checksum(buf[:], crc64.MakeTable(crc64.ECMA))
	id := fmt.Sprintf("%X", sum)

	inst := &instance{
		Instance:         trace.NewInstance(dim, sceneFn()),
		sceneName:        sceneName,
		created:          time.Now(),
		dim:              dim,
		requestedWorkers: 0,
	}
	c.instances[id] = inst
	return id, nil
}

func (c *controller) setWorkers(renderID string, workers int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	inst, ok := c.instances[renderID]
	if !ok {
		return fmt.Errorf("unknown render id: %v", renderID)
	}
	inst.requestedWorkers = workers
	inst.SetWorkers(workers)
	return nil
}

func (c *controller) getImage(renderID string) (image.Image, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	inst, ok := c.instances[renderID]
	if !ok {
		return nil, fmt.Errorf("unknown render id: %v", renderID)
	}
	return inst.Image(), nil
}
