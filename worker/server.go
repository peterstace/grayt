package worker

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func NewServer() *Server {
	return new(Server)
}

type Server struct {
	mu      sync.Mutex
	working bool
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/trace" {
		s.handleTrace(w, req)
		return
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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

	// TODO Get the scene if it hasn't been cached locally.

	// TODO Trace the image and return the data.

	time.Sleep(time.Second)
	fmt.Fprintf(w, "DATA GETS SEND HERE")

}
