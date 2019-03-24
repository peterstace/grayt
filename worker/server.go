package worker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	working bool

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
	if s.working {
		http.Error(w, "already working", http.StatusTooManyRequests)
		s.mu.Unlock()
		return
	}
	s.working = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.working = false
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

	// TODO Trace the image and return the data.

	time.Sleep(time.Second)
	fmt.Fprintf(w, "DATA GETS SEND HERE")
}
