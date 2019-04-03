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
	Workers     int
	Completed   int
	Passes      int
	TraceRateHz int
}

type Instance struct {
	accel accelerationStructure
	cam   camera

	accum *accumulator

	reqWorkersCond   *sync.Cond
	requestedWorkers int

	actualWorkers int64
	completed     int64
	traceRate     int64
	workIndex     int64

	workerWG sync.WaitGroup
}

func NewInstance(dim xmath.Dimensions, scn scene.Scene) *Instance {
	cam, objs := buildScene(scn)
	inst := &Instance{
		accel:          newGrid(4, objs),
		cam:            cam,
		accum:          newAccumulator(dim),
		reqWorkersCond: sync.NewCond(new(sync.Mutex)),
	}
	go inst.dispatchWork()
	go inst.monitorTraceRate()
	return inst
}

func (in *Instance) SetWorkers(workers int) {
	in.reqWorkersCond.L.Lock()
	in.requestedWorkers = workers
	in.reqWorkersCond.L.Unlock()
	in.reqWorkersCond.Broadcast()
}

func (in *Instance) dispatchWork() {
	for {
		cnd := in.reqWorkersCond
		cnd.L.Lock()
		for in.requestedWorkers == 0 {
			atomic.StoreInt64(&in.actualWorkers, 0)
			cnd.Wait()
		}

		atomic.StoreInt64(&in.workIndex, 0)
		atomic.StoreInt64(&in.actualWorkers, int64(in.requestedWorkers))
		for i := 0; i < in.requestedWorkers; i++ {
			in.workerWG.Add(1)
			go in.work()
		}
		cnd.L.Unlock()

		in.workerWG.Wait()
		in.accum.merge(1)
	}
}

func (in *Instance) GetStats() Stats {
	return Stats{
		Workers:     int(atomic.LoadInt64(&in.actualWorkers)),
		Completed:   int(atomic.LoadInt64(&in.completed)),
		Passes:      in.accum.getPasses(),
		TraceRateHz: int(atomic.LoadInt64(&in.traceRate)),
	}
}

func (in *Instance) work() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tr := newTracer(in.accel, rng)
	wide := in.accum.dim.Wide
	high := in.accum.dim.High
	pxPitch := 2.0 / float64(wide)
	for {
		idx := int(atomic.AddInt64(&in.workIndex, 1))
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
	in.workerWG.Done()
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
