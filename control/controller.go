package control

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

func New() *Controller {
	return &Controller{
		instances: make(map[string]*instance),
	}
}

type Controller struct {
	mu        sync.Mutex
	instances map[string]*instance
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
			ID:               id,
			SceneName:        inst.sceneName,
			Created:          inst.created,
			Dimensions:       inst.accum.dim,
			Passes:           inst.getPasses(),
			Completed:        inst.getCompleted(),
			TraceRateHz:      float64(inst.getTraceRateHz()),
			RequestedWorkers: inst.requestedWorkers,
			ActualWorkers:    inst.getWorkers(),
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

	scn := trace.BuildScene(sceneFn())
	accel := trace.NewGrid(4, scn.Objects)

	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(time.Now().Unix()))
	sum := crc64.Checksum(buf[:], crc64.MakeTable(crc64.ECMA))
	id := fmt.Sprintf("%X", sum)

	inst := newInstance(sceneName, dim, accel, scn.Camera)
	go inst.dispatchWork()
	c.instances[id] = inst
	return id, nil
}

func (c *Controller) SetWorkers(renderID string, workers int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	inst, ok := c.instances[renderID]
	if !ok {
		return fmt.Errorf("unknown render id: %v", renderID)
	}
	inst.setWorkers(workers)
	return nil
}

func (c *Controller) GetImage(renderID string) (image.Image, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	inst, ok := c.instances[renderID]
	if !ok {
		return nil, fmt.Errorf("unknown render id: %v", renderID)
	}
	return inst.accum.toImage(1.0), nil
}
