package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
)

func NewServer(scenelibAddr, workerAddr, assetsDir string) *Server {
	return &Server{
		scenelibAddr: scenelibAddr,
		workerAddr:   workerAddr,
		assets:       http.FileServer(http.Dir(assetsDir)),
		renders:      map[string]render{},
	}
}

type Server struct {
	scenelibAddr string
	workerAddr   string
	assets       http.Handler

	mu      sync.Mutex
	renders map[string]render
}

type render struct {
	Scene  string
	PxWide int
	PxHigh int
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/scenes" {
		s.handleGetScenes(w, req)
		return
	}
	if req.URL.Path == "/renders" {
		switch req.Method {
		case http.MethodGet:
			s.handleGetRenders(w, req)
		case http.MethodPost:
			s.handlePostRenders(w, req)
		default:
			http.Error(w, "method must be GET or POST", http.StatusMethodNotAllowed)
		}
		return
	}
	s.assets.ServeHTTP(w, req)
}

func (s *Server) handleGetScenes(w http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("http://" + s.scenelibAddr + "/scenes")
	if err != nil {
		http.Error(w,
			"fetching scene list: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintf(w, "fetching scene list: ")
		io.Copy(w, resp.Body)
		return
	}
	io.Copy(w, resp.Body)
}

func (s *Server) handleGetRenders(w http.ResponseWriter, req *http.Request) {
	type resource struct {
		Scene            string `json:"scene"`
		PxWide           int    `json:"px_wide"`
		PxHigh           int    `json:"px_high"`
		Passes           int    `json:"passes"`
		Completed        string `json:"completed"`
		TraceRate        string `json:"trace_rate"`
		ID               string `json:"uuid"`
		RequestedWorkers int    `json:"requested_workers"`
		ActualWorkers    int    `json:"actual_workers"`
	}
	resources := []resource{} // init as empty array because it marshals to json

	s.mu.Lock()
	for id, r := range s.renders {
		resources = append(resources, resource{
			Scene:            r.Scene,
			PxWide:           r.PxWide,
			PxHigh:           r.PxHigh,
			Passes:           0,
			Completed:        "0",
			TraceRate:        "0",
			ID:               id,
			RequestedWorkers: 0,
			ActualWorkers:    0,
		})
	}
	s.mu.Unlock()
	// TODO: order resources by created at time

	if err := json.NewEncoder(w).Encode(resources); err != nil {
		http.Error(w, "encoding renders: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handlePostRenders(w http.ResponseWriter, req *http.Request) {
	var form struct {
		Scene  string `json:"scene"`
		PxWide int    `json:"px_wide"`
		PxHigh int    `json:"px_high"`
	}
	if err := json.NewDecoder(req.Body).Decode(&form); err != nil {
		http.Error(w, "decoding form: "+err.Error(), http.StatusBadRequest)
		return
	}
	if form.PxWide == 0 || form.PxHigh == 0 {
		http.Error(w, "px_wide or px_high not set", http.StatusBadRequest)
		return
	}

	// TODO: Check if scene exists. Should call scenelib service.

	id := fmt.Sprintf("%X", rand.Uint64())

	s.mu.Lock()
	s.renders[id] = render{
		Scene:  form.Scene,
		PxWide: form.PxWide,
		PxHigh: form.PxHigh,
	}
	s.mu.Unlock()

	fmt.Fprintf(w, `{"uuid":%q}`, id)
}
