package scenelib

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/peterstace/grayt/scene"
)

func NewServer() *Server {
	return &Server{
		sceneCache: make(map[string]scene.Scene),
		registry:   make(map[string]func() scene.Scene),
	}
}

type Server struct {
	mu         sync.Mutex
	sceneCache map[string]scene.Scene

	registry map[string]func() scene.Scene
}

func (s *Server) Register(
	name string,
	sceneFn func() scene.Scene,
) {
	s.registry[name] = sceneFn
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/scene" {
		s.handleScene(w, req)
		return
	}
	if req.URL.Path == "/scenes" {
		s.handleScenes(w, req)
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

func (s *Server) handleScenes(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method must be GET", http.StatusBadRequest)
		return
	}
	type scene struct {
		Code string `json:"code"`
	}
	var scenes []scene
	for code := range s.registry {
		scenes = append(scenes, scene{code})
	}
	if err := json.NewEncoder(w).Encode(scenes); err != nil {
		http.Error(w, "encoding scenes: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
