package worker

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/peterstace/grayt/protocol"
	"github.com/peterstace/grayt/trace"
)

func NewServer(scenelibAddr string) *Server {
	return &Server{
		scenelibAddr: scenelibAddr,
		scenes:       map[string]trace.Scene{},
	}
}

type Server struct {
	scenelibAddr string

	mu      sync.Mutex
	serving int

	scenes map[string]trace.Scene
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
	s.serveLayer(w, sceneName, pxWide, pxHigh)
}

func (s *Server) serveLayer(
	w http.ResponseWriter, sceneName string, pxWide, pxHigh int,
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

	scene, ok := s.scenes[sceneName]
	if !ok {
		log.Printf("scene %q not cached, fetching from scenelib", sceneName)
		resp, err := http.Get(
			s.scenelibAddr + "/scene?name=" + url.PathEscape(sceneName),
		)
		if err != nil {
			http.Error(w,
				"fetching scene: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(resp.StatusCode)
			fmt.Fprintf(w, "fetching scene: ")
			io.Copy(w, resp.Body)
			return
		}
		var sceneProto protocol.Scene
		if err := json.NewDecoder(resp.Body).Decode(&sceneProto); err != nil {
			http.Error(w,
				fmt.Sprintf("decoding scene: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
		scene = trace.BuildScene(sceneProto)
		s.scenes[sceneName] = scene
	}

	traceLayer(w, pxWide, pxHigh, scene)
}

func traceLayer(w io.Writer, pxWide, pxHigh int, scene trace.Scene) {
	// TODO: a lot of this can be cached
	accel := trace.NewGrid(4, scene.Objects)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tr := trace.NewTracer(accel, rng)
	pxPitch := 2.0 / float64(pxWide)
	for pxY := 0; pxY < pxHigh; pxY++ {
		for pxX := 0; pxX < pxWide; pxX++ {
			x := (float64(pxX-pxWide/2) + rng.Float64()) * pxPitch
			y := (float64(pxY-pxHigh/2) + rng.Float64()) * pxPitch * -1.0
			cr := scene.Camera.MakeRay(x, y, rng)
			cr.Dir = cr.Dir.Unit()
			c := tr.TracePath(cr)
			binary.Write(w, binary.BigEndian, c.R)
			binary.Write(w, binary.BigEndian, c.G)
			binary.Write(w, binary.BigEndian, c.B)
		}
	}
}
