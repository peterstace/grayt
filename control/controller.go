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
		renders: make(map[string]Render),
	}
}

type Controller struct {
	mu      sync.Mutex
	renders map[string]Render
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
	for _, ren := range c.renders {
		renders = append(renders, ren)
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

	c.renders[id] = Render{
		ID:         id,
		SceneName:  sceneName,
		Dimensions: dim,
		Created:    time.Now(),
	}
	return id, nil
}

func (c *Controller) SetWorkers(renderID string, workers int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	ren, ok := c.renders[renderID]
	if !ok {
		return fmt.Errorf("unknown render id: %v", renderID)
	}
	ren.RequestedWorkers = workers
	c.renders[renderID] = ren
	return nil
}
