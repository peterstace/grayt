package api

import (
	"encoding/json"
	"net/http"

	"github.com/peterstace/grayt/scene/library"
)

func (s *Server) handleGetScenes(w http.ResponseWriter) {
	type scn struct {
		Code string `json:"code"`
	}
	var scns []scn
	for _, name := range library.Listing() {
		scns = append(scns, scn{name})
	}
	if err := json.NewEncoder(w).Encode(scns); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetRenders(w http.ResponseWriter) {
	// TODO
}

func (s *Server) handlePostRenders(w http.ResponseWriter, req *http.Request) {
	// TODO
}

func (s *Server) handlePutWorkers(w http.ResponseWriter, req *http.Request, id string) {
	// TODO
}

func (s *Server) handleGetImage(w http.ResponseWriter, id string) {
	// TODO
}
