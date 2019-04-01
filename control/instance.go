package control

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/peterstace/grayt/trace"
	"github.com/peterstace/grayt/xmath"
)

// TODO: this functionality should live in the trace package (maybe called 'Tracer')

type instance struct {
	sceneName string
	created   time.Time
	accel     trace.AccelerationStructure
	cam       trace.Camera

	accum *accumulator

	reqWorkersCond   *sync.Cond
	requestedWorkers int

	actualWorkers int64
	completed     int64
	traceRate     int64
	workIndex     int64

	workerWG sync.WaitGroup
}

func newInstance(sceneName string, dim xmath.Dimensions, accel trace.AccelerationStructure, cam trace.Camera) *instance {
	inst := &instance{
		sceneName:      sceneName,
		created:        time.Now(),
		accel:          accel,
		cam:            cam,
		accum:          newAccumulator(dim),
		reqWorkersCond: sync.NewCond(new(sync.Mutex)),
	}
	go inst.monitorTraceRate()
	return inst
}

func (in *instance) setWorkers(workers int) {
	in.reqWorkersCond.L.Lock()
	in.requestedWorkers = workers
	in.reqWorkersCond.L.Unlock()
	in.reqWorkersCond.Broadcast()
}

func (in *instance) dispatchWork() {
	for {
		cnd := in.reqWorkersCond
		cnd.L.Lock()
		for in.requestedWorkers == 0 {
			atomic.StoreInt64(&in.actualWorkers, 0)
			cnd.Wait()
		}

		// TODO: loop over number of workers and launch a goroutine for each.
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

func (in *instance) getWorkers() int {
	return int(atomic.LoadInt64(&in.actualWorkers))
}

func (in *instance) getCompleted() int {
	return int(atomic.LoadInt64(&in.completed))
}

func (in *instance) getPasses() int {
	return in.accum.getPasses()
}

func (in *instance) work() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tr := trace.NewTracer(in.accel, rng)
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
		y := (float64(pxY-wide/2) + rng.Float64()) * pxPitch * -1.0
		cr := in.cam.MakeRay(x, y, rng)
		cr.Dir = cr.Dir.Unit()
		c := tr.TracePath(cr)
		in.accum.set(pxX, pxY, c)
		atomic.AddInt64(&in.completed, 1)
	}
	in.workerWG.Done()
}

func (in *instance) getTraceRateHz() int {
	return int(atomic.LoadInt64(&in.traceRate))
}

func (in *instance) monitorTraceRate() {
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
