package api

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"image/png"
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
		s.mu.Lock()
		ren, ok := s.renders[id]
		if !ok {
			http.Error(w, "unknown render id", http.StatusBadRequest)
			s.mu.Unlock()
			return
		}
		s.mu.Unlock()
		switch parts[1] {
		case "workers":
			if req.Method != http.MethodPut {
				http.Error(w, "method must be PUT", http.StatusMethodNotAllowed)
				return
			}
			s.handlePutWorkers(w, req, ren)
		case "image":
			if req.Method != http.MethodGet {
				http.Error(w, "method must be GET", http.StatusMethodNotAllowed)
				return
			}
			s.handleGetImage(w, req, ren)
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
		passes := int(r.acc.getPasses())
		r.cnd.L.Lock()
		resources = append(resources, resource{
			Scene:            r.scene,
			PxWide:           r.pxWide,
			PxHigh:           r.pxHigh,
			Passes:           passes,
			Completed:        displayFloat64(float64(passes * r.pxHigh * r.pxHigh)),
			TraceRate:        displayFloat64(r.monitor.rateHz()),
			ID:               id,
			RequestedWorkers: r.desiredWorkers,
			ActualWorkers:    r.actualWorkers,
		})
		r.cnd.L.Unlock()
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

	newRender := render{
		scene:      form.Scene,
		pxWide:     form.PxWide,
		pxHigh:     form.PxHigh,
		created:    time.Now(),
		cnd:        sync.NewCond(new(sync.Mutex)),
		acc:        newAccumulator(form.PxWide, form.PxHigh),
		monitor:    rateMonitor{points: list.New()},
		workerAddr: s.workerAddr,
	}

	id := fmt.Sprintf("%X", rand.Uint64())
	s.mu.Lock()
	s.renders[id] = &newRender
	s.mu.Unlock()

	go newRender.orchestrateWork()

	fmt.Fprintf(w, `{"uuid":%q}`, id)
}

func (s *Server) handlePutWorkers(
	w http.ResponseWriter, req *http.Request, ren *render,
) {
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

func (s *Server) handleGetImage(
	w http.ResponseWriter, req *http.Request, ren *render,
) {
	var buf bytes.Buffer
	img := ren.acc.toImage(1.0)
	if err := png.Encode(&buf, img); err != nil {
		http.Error(w,
			"could not encode image: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	io.Copy(w, &buf)
}

func displayFloat64(f float64) string {
	var thousands int
	for f >= 1000 {
		f /= 1000
		thousands++
	}
	var body string
	switch {
	case f < 10:
		body = fmt.Sprintf("%.3f", f)
	case f < 100:
		body = fmt.Sprintf("%.2f", f)
	case f < 1000:
		body = fmt.Sprintf("%.1f", f)
	default:
		panic(f)
	}
	return fmt.Sprintf("%se%d", body, thousands*3)
}
