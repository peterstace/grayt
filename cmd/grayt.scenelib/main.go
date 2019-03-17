package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/peterstace/grayt/examples/cornellbox/classic"
	"github.com/peterstace/grayt/examples/cornellbox/reflections"
	"github.com/peterstace/grayt/examples/cornellbox/spheretree"
	"github.com/peterstace/grayt/examples/cornellbox/splitbox"
	"github.com/peterstace/grayt/grayt"
)

const listenAddrEnv = "GRAYT_SCENELIB_LISTEN_ADDR"

func main() {
	listenAddr, ok := os.LookupEnv(listenAddrEnv)
	if !ok {
		log.Fatalf("%s not set", listenAddrEnv)
	}

	s := Server{
		sceneCache: make(map[string]grayt.Scene),
		registry:   make(map[string]func() grayt.Scene),
	}
	s.Register("cornellbox_classic", classic.CameraFn(), classic.ObjectsFn)
	s.Register("cornellbox_reflections", reflections.CameraFn(), reflections.ObjectsFn)
	s.Register("spheretree", spheretree.CameraFn(), spheretree.ObjectsFn)
	s.Register("splitbox", splitbox.CameraFn(), splitbox.ObjectsFn)

	http.HandleFunc("/scene", s.HandleScene)
	log.Printf("serving on %v", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

type Server struct {
	mu         sync.Mutex
	sceneCache map[string]grayt.Scene

	registry map[string]func() grayt.Scene
}

func (s *Server) Register(
	name string,
	cam grayt.CameraBlueprint,
	objs func() grayt.ObjectList,
) {
	s.registry[name] = func() grayt.Scene {
		return grayt.Scene{
			Camera:  cam,
			Objects: objs(),
		}
	}
}

func (s *Server) HandleScene(w http.ResponseWriter, req *http.Request) {
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
