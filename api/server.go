package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func NewServer(scenelibAddr, workerAddr, assetsDir string) *Server {
	return &Server{
		scenelibAddr: scenelibAddr,
		workerAddr:   workerAddr,
		assets:       http.FileServer(http.Dir(assetsDir)),
		renders:      map[string]*render{},
	}
}

type Server struct {
	scenelibAddr string
	workerAddr   string
	assets       http.Handler

	mu      sync.Mutex
	renders map[string]*render
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

	if strings.HasPrefix(req.URL.Path, "/renders/") {
		rest := strings.TrimPrefix(req.URL.Path, "/renders/")
		parts := strings.Split(rest, "/")
		if len(parts) != 2 {
			http.NotFound(w, req)
			return
		}
		id := parts[0]
		switch parts[1] {
		case "workers":
			if req.Method != http.MethodPut {
				http.Error(w, "method must be PUT", http.StatusMethodNotAllowed)
				return
			}
			s.handlePutWorkers(w, req, id)
		case "image":
			if req.Method != http.MethodGet {
				http.Error(w, "method must be GET", http.StatusMethodNotAllowed)
				return
			}
			s.handleGetImage(w, req)
		default:
			http.NotFound(w, req)
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
		r.cnd.L.Lock()
		requested := r.desiredWorkers
		r.cnd.L.Unlock()
		resources = append(resources, resource{
			Scene:            r.scene,
			PxWide:           r.pxWide,
			PxHigh:           r.pxHigh,
			Passes:           0,
			Completed:        "0",
			TraceRate:        "0",
			ID:               id,
			RequestedWorkers: requested,
			ActualWorkers:    0,
		})
	}
	sort.Slice(resources, func(i, j int) bool {
		t1 := s.renders[resources[i].ID].created
		t2 := s.renders[resources[j].ID].created
		return t1.Before(t2)
	})
	s.mu.Unlock()

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

	newRender := render{
		scene:   form.Scene,
		pxWide:  form.PxWide,
		pxHigh:  form.PxHigh,
		created: time.Now(),
		cnd:     sync.NewCond(new(sync.Mutex)),
	}

	s.mu.Lock()
	s.renders[id] = &newRender
	s.mu.Unlock()

	go newRender.work()

	fmt.Fprintf(w, `{"uuid":%q}`, id)
}

func (s *Server) handlePutWorkers(
	w http.ResponseWriter, req *http.Request, id string,
) {
	s.mu.Lock()
	ren, ok := s.renders[id]
	if !ok {
		http.Error(w, "unknown render id", http.StatusBadRequest)
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "could not read body", http.StatusInternalServerError)
		return
	}
	workers, err := strconv.Atoi(string(buf))
	if err != nil {
		http.Error(w, "could not parse worker count", http.StatusBadRequest)
		return
	}
	if workers < 0 {
		http.Error(w, "workers must be non-negative", http.StatusBadRequest)
		return
	}

	ren.cnd.L.Lock()
	ren.desiredWorkers = workers
	ren.cnd.L.Unlock()
	ren.cnd.Broadcast()

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetImage(w http.ResponseWriter, req *http.Request) {
	// TODO
	w.WriteHeader(http.StatusOK)
}
