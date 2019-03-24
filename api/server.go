package api

import (
	"fmt"
	"io"
	"net/http"
)

func NewServer(scenelibAddr, workerAddr, assetsDir string) *Server {
	return &Server{
		scenelibAddr: scenelibAddr,
		workerAddr:   workerAddr,
		assets:       http.FileServer(http.Dir(assetsDir)),
	}
}

type Server struct {
	scenelibAddr string
	workerAddr   string
	assets       http.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/scenes" {
		s.handleGetScenes(w, req)
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
