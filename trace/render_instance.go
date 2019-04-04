package trace

import (
	"image"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/peterstace/grayt/scene"
	"github.com/peterstace/grayt/xmath"
)

type Stats struct {
	LoadState   string
	Workers     int
	Completed   int
	Passes      int
	TraceRateHz int
}

// States: unloaded, loading, loaded

type loadState int

const (
	unloaded loadState = iota
	loading
	loaded
)

type Instance struct {
	// Access controlled by cond variable
	cond             *sync.Cond
	requestedWorkers int
	loadState        loadState
	sceneFn          func() scene.Scene
	accel            accelerationStructure
	cam              camera

	// Access self controlled
	accum *accumulator

	// Atomic access
	actualWorkers int64
	completed     int64
	traceRate     int64
}

func NewInstance(dim xmath.Dimensions, sceneFn func() scene.Scene) *Instance {
	inst := &Instance{
		sceneFn: sceneFn,
		accum:   newAccumulator(dim),
		cond:    sync.NewCond(new(sync.Mutex)),
	}
	go inst.loadScene()
	go inst.dispatchWork()
	go inst.monitorTraceRate()
	return inst
}

func (in *Instance) loadScene() {
	in.cond.L.Lock()
	for in.requestedWorkers == 0 || in.loadState != unloaded {
		in.cond.Wait()
	}
	in.loadState = loading
	in.cond.L.Unlock()

	cam, objs := buildScene(in.sceneFn())
	in.cam = cam
	in.accel = newGrid(4, objs)

	in.cond.L.Lock()
	in.loadState = loaded
	in.cond.Broadcast()
	in.cond.L.Unlock()
}

func (in *Instance) SetWorkers(workers int) {
	in.cond.L.Lock()
	in.requestedWorkers = workers
	in.cond.L.Unlock()
	in.cond.Broadcast()
}

func (in *Instance) dispatchWork() {
	for {
		in.cond.L.Lock()
		for in.requestedWorkers == 0 || in.loadState != loaded {
			atomic.StoreInt64(&in.actualWorkers, 0)
			in.cond.Wait()
		}

		var ctx workContext
		atomic.StoreInt64(&in.actualWorkers, int64(in.requestedWorkers))
		ctx.wg.Add(in.requestedWorkers)
		for i := 0; i < in.requestedWorkers; i++ {
			go in.work(&ctx)
		}
		in.cond.L.Unlock()

		ctx.wg.Wait()
		in.accum.merge(1)
	}
}

func (in *Instance) GetStats() Stats {
	in.cond.L.Lock()
	loadState := map[loadState]string{
		unloaded: "unloaded",
		loading:  "loading",
		loaded:   "loaded",
	}[in.loadState]
	in.cond.L.Unlock()

	return Stats{
		LoadState:   loadState,
		Workers:     int(atomic.LoadInt64(&in.actualWorkers)),
		Completed:   int(atomic.LoadInt64(&in.completed)),
		Passes:      in.accum.getPasses(),
		TraceRateHz: int(atomic.LoadInt64(&in.traceRate)),
	}
}

type workContext struct {
	idx int64
	wg  sync.WaitGroup
}

func (in *Instance) work(ctx *workContext) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tr := newTracer(in.accel, rng)
	wide := in.accum.dim.Wide
	high := in.accum.dim.High
	pxPitch := 2.0 / float64(wide)
	for {
		idx := int(atomic.AddInt64(&ctx.idx, 1))
		if idx >= wide*high {
			break
		}
		pxY := idx / wide
		pxX := idx % wide
		x := (float64(pxX-wide/2) + rng.Float64()) * pxPitch
		y := (float64(pxY-high/2) + rng.Float64()) * pxPitch * -1.0
		cr := in.cam.makeRay(x, y, rng)
		cr.Dir = cr.Dir.Unit()
		c := tr.tracePath(cr)
		in.accum.set(pxX, pxY, c)
		atomic.AddInt64(&in.completed, 1)
	}
	ctx.wg.Done()
}

func (in *Instance) monitorTraceRate() {
	const samplePeriod = 5 * time.Second
	var lastCompleted int64
	ticker := time.NewTicker(samplePeriod)
	for {
		<-ticker.C
		completed := atomic.LoadInt64(&in.completed)
		if lastCompleted != 0 {
			sample := (completed - lastCompleted) * int64(time.Second) / int64(samplePeriod)
			atomic.StoreInt64(&in.traceRate, sample)
		}
		lastCompleted = completed
	}
}

func (in *Instance) Image() image.Image {
	return in.accum.toImage(1.0)
}
