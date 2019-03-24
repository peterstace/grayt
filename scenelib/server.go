package scenelib

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/peterstace/grayt/protocol"
)

func NewServer() *Server {
	return &Server{
		sceneCache: make(map[string]protocol.Scene),
		registry:   make(map[string]func() protocol.Scene),
	}
}

type Server struct {
	mu         sync.Mutex
	sceneCache map[string]protocol.Scene

	registry map[string]func() protocol.Scene
}

func (s *Server) Register(
	name string,
	sceneFn func() protocol.Scene,
) {
	s.registry[name] = sceneFn
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/scene" {
		s.handleScene(w, req)
		return
	}
	http.NotFound(w, req)
}

func (s *Server) handleScene(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method must be GET", http.StatusBadRequest)
		return
	}
	name := req.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "query parameter 'name' not set", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	scene, ok := s.sceneCache[name]
	if ok {
		s.mu.Unlock()
	} else {
		sceneFn, ok := s.registry[name]
		if !ok {
			http.Error(w, "unknown scene name", http.StatusBadRequest)
			s.mu.Unlock()
			return
		}
		scene = sceneFn()
		s.sceneCache[name] = scene
		s.mu.Unlock()
	}

	if err := json.NewEncoder(w).Encode(scene); err != nil {
		http.Error(w, "couldn't write scene: "+err.Error(), http.StatusInternalServerError)
		return
	}
}