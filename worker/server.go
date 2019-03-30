package worker

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/scene/library"
	"github.com/peterstace/grayt/trace"
)

func NewServer() *Server {
	return &Server{
		scenes: map[string]tracer{},
	}
}

type Server struct {
	mu      sync.Mutex
	serving int

	scenes map[string]tracer
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/trace" {
		s.handleTrace(w, req)
		return
	}
	http.NotFound(w, req)
}

func (s *Server) handleTrace(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	params := req.URL.Query()
	sceneName := params.Get("scene_name")
	if sceneName == "" {
		http.Error(w, "scene_name query param not set", http.StatusBadRequest)
		return
	}
	pxWideStr := params.Get("px_wide")
	if pxWideStr == "" {
		http.Error(w, "px_wide query param not set", http.StatusBadRequest)
		return
	}
	pxWide, err := strconv.Atoi(pxWideStr)
	if err != nil {
		http.Error(w, "couldn't convert px_wide to int", http.StatusBadRequest)
		return
	}
	pxHighStr := params.Get("px_high")
	if pxHighStr == "" {
		http.Error(w, "px_high query param not set", http.StatusBadRequest)
		return
	}
	pxHigh, err := strconv.Atoi(pxHighStr)
	if err != nil {
		http.Error(w, "couldn't convert px_high to int", http.StatusBadRequest)
		return
	}
	depthStr := params.Get("depth")
	if depthStr == "" {
		http.Error(w, "depth query param not set", http.StatusBadRequest)
		return
	}
	depth, err := strconv.Atoi(depthStr)
	if err != nil {
		http.Error(w, "couldn't convert depth int", http.StatusBadRequest)
		return
	}
	s.serveLayer(w, sceneName, pxWide, pxHigh, depth)
}

func (s *Server) serveLayer(
	w http.ResponseWriter, sceneName string, pxWide, pxHigh, depth int,
) {
	s.mu.Lock()
	if s.serving >= 4 {
		http.Error(w, "already working", http.StatusTooManyRequests)
		s.mu.Unlock()
		return
	}
	s.serving++
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.serving--
		s.mu.Unlock()
	}()

	tr, ok := s.scenes[sceneName]
	if !ok {
		sceneFn, ok := library.Lookup(sceneName)
		if !ok {
			http.Error(w, fmt.Sprintf("unknown scene: %v", sceneName), http.StatusBadRequest)
			return
		}
		scene := trace.BuildScene(sceneFn())
		accel := trace.NewGrid(4, scene.Objects)
		tr = tracer{scene.Camera, accel}
		s.scenes[sceneName] = tr
	}
	tr.traceLayer(w, pxWide, pxHigh, depth)
}

type tracer struct {
	camera trace.Camera
	accel  trace.AccelerationStructure
}

func (t *tracer) traceLayer(w io.Writer, pxWide, pxHigh, depth int) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tr := trace.NewTracer(t.accel, rng)
	pxPitch := 2.0 / float64(pxWide)
	for pxY := 0; pxY < pxHigh; pxY++ {
		for pxX := 0; pxX < pxWide; pxX++ {
			var c colour.Colour
			for i := 0; i < depth; i++ {
				x := (float64(pxX-pxWide/2) + rng.Float64()) * pxPitch
				y := (float64(pxY-pxHigh/2) + rng.Float64()) * pxPitch * -1.0
				cr := t.camera.MakeRay(x, y, rng)
				cr.Dir = cr.Dir.Unit()
				c = c.Add(tr.TracePath(cr))
			}
			binary.Write(w, binary.BigEndian, c)
		}
	}
}
