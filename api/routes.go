package api

import (
	"fmt"
	"net/http"
	"strings"
)

func NewServer(assetsDir string) *Server {
	return &Server{
		assets: http.FileServer(http.Dir(assetsDir)),
		ctrl:   newController(),
	}
}

type Server struct {
	assets http.Handler
	ctrl   *controller
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logRequests(http.HandlerFunc(s.routeRoot)).ServeHTTP(w, req)
}

func (s *Server) routeRoot(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/scenes" {
		s.routeScenes(w, req)
		return
	}
	if strings.HasPrefix(req.URL.Path, "/renders") {
		s.routeRenders(w, req)
		return
	}
	s.assets.ServeHTTP(w, req)
}

func (s *Server) routeScenes(w http.ResponseWriter, req *http.Request) {
	if !methodAllowed(w, req, http.MethodGet) {
		return
	}
	s.handleGetScenes(w)
}

func (s *Server) routeRenders(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/renders" {
		if !methodAllowed(w, req, http.MethodGet, http.MethodPost) {
			return
		}
		if req.Method == http.MethodGet {
			s.handleGetRenders(w)
		} else {
			s.handlePostRenders(w, req)
		}
		return
	}
	if strings.HasPrefix(req.URL.Path, "/renders/") {
		rest := strings.TrimPrefix(req.URL.Path, "/renders/")
		parts := strings.Split(rest, "/")
		if len(parts) == 2 {
			id := parts[0]
			switch parts[1] {
			case "workers":
				s.routeWorkers(w, req, id)
				return
			case "image":
				s.routeImage(w, req, id)
				return
			}
		}
	}
	http.NotFound(w, req)
}

func (s *Server) routeWorkers(w http.ResponseWriter, req *http.Request, id string) {
	if !methodAllowed(w, req, http.MethodPut) {
		return
	}
	s.handlePutWorkers(w, req, id)
}

func (s *Server) routeImage(w http.ResponseWriter, req *http.Request, id string) {
	if !methodAllowed(w, req, http.MethodGet) {
		return
	}
	s.handleGetImage(w, id)
}

func methodAllowed(w http.ResponseWriter, req *http.Request, allowed ...string) bool {
	for _, allow := range allowed {
		if req.Method == allow {
			return true
		}
	}
	msg := fmt.Sprintf(
		"%s: must be %s",
		http.StatusText(http.StatusMethodNotAllowed),
		strings.Join(allowed, ","),
	)
	http.Error(w, msg, http.StatusMethodNotAllowed)
	return false
}
