package api

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"image"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/peterstace/grayt/scene/library"
	"github.com/peterstace/grayt/trace"
	"github.com/peterstace/grayt/xmath"
)

func newController(dataDir string) *controller {
	return &controller{
		dataDir:   dataDir,
		instances: make(map[string]*instance),
	}
}

type controller struct {
	dataDir string

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
	LoadState        string    `json:"load_state"`
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
			LoadState:        stats.LoadState,
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

type save struct {
	SceneName string           `json:"scene_name"`
	Created   time.Time        `json:"created"`
	Dim       xmath.Dimensions `json:"dim"`
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

	accumFilename := filepath.Join(c.dataDir, id+".data")
	metadataFilename := filepath.Join(c.dataDir, id+".json")

	now := time.Now()
	metadataBuf, err := json.Marshal(save{
		SceneName: sceneName,
		Created:   now,
		Dim:       dim,
	})
	if err != nil {
		return "", fmt.Errorf("marshalling metadata: %v", err)
	}
	if err := ioutil.WriteFile(metadataFilename, metadataBuf, 0775); err != nil {
		return "", fmt.Errorf("writing metadata: %v", err)
	}

	inst := &instance{
		Instance:         trace.NewInstance(dim, sceneFn, accumFilename),
		sceneName:        sceneName,
		created:          now,
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
