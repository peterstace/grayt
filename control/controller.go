package control

import (
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"sync"
	"time"

	"github.com/peterstace/grayt/scene/library"
	"github.com/peterstace/grayt/xmath"
)

func New() *Controller {
	return &Controller{
		instances: make(map[string]*instance),
	}
}

type Controller struct {
	mu        sync.Mutex
	instances map[string]*instance
}

type instance struct {
	sceneName        string
	created          time.Time
	requestedWorkers int
	accum            *accumulator
}

type Render struct {
	ID               string
	SceneName        string
	Created          time.Time
	Dimensions       xmath.Dimensions
	Passes           int
	Completed        int
	TraceRateHz      float64
	RequestedWorkers int
	ActualWorkers    int
}

func (c *Controller) GetRenders() []Render {
	c.mu.Lock()
	defer c.mu.Unlock()
	var renders []Render
	for id, inst := range c.instances {
		renders = append(renders, Render{
			ID:        id,
			SceneName: inst.sceneName,
			Created:   inst.created,
			Dimensions: xmath.Dimensions{
				Wide: inst.accum.wide,
				High: inst.accum.high,
			},
			Passes:           0,
			Completed:        0,
			TraceRateHz:      0.0,
			RequestedWorkers: inst.requestedWorkers,
			ActualWorkers:    0,
		})
	}
	return renders
}

func (c *Controller) NewRender(sceneName string, dim xmath.Dimensions) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	sceneFn, ok := library.Lookup(sceneName)
	if !ok {
		return "", fmt.Errorf("unknown scene name: %v", sceneName)
	}

	_ = sceneFn // TODO: should build out scene and store

	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(time.Now().Unix()))
	sum := crc64.Checksum(buf[:], crc64.MakeTable(crc64.ECMA))
	id := fmt.Sprintf("%X", sum)

	c.instances[id] = &instance{
		sceneName: sceneName,
		created:   time.Now(),
		accum:     newAccumulator(dim.Wide, dim.High),
	}
	return id, nil
}

func (c *Controller) SetWorkers(renderID string, workers int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	inst, ok := c.instances[renderID]
	if !ok {
		return fmt.Errorf("unknown render id: %v", renderID)
	}
	inst.requestedWorkers = workers
	return nil
}
