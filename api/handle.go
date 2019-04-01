package api

import (
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	"github.com/peterstace/grayt/scene/library"
	"github.com/peterstace/grayt/xmath"
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

	renders := s.ctrl.GetRenders()
	sort.Slice(renders, func(i, j int) bool {
		return renders[i].Created.Before(renders[j].Created)
	})
	for _, r := range renders {
		resources = append(resources, resource{
			Scene:            r.SceneName,
			PxWide:           r.Dimensions.Wide,
			PxHigh:           r.Dimensions.High,
			Passes:           r.Passes,
			Completed:        displayFloat64(float64(r.Completed)),
			TraceRate:        displayFloat64(r.TraceRateHz) + " Hz",
			ID:               r.ID,
			RequestedWorkers: r.RequestedWorkers,
			ActualWorkers:    r.ActualWorkers,
		})
	}

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

	id, err := s.ctrl.NewRender(form.Scene, xmath.Dimensions{form.PxWide, form.PxHigh})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, `{"uuid":%q}`, id)
}

func (s *Server) handlePutWorkers(w http.ResponseWriter, req *http.Request, id string) {
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

	if err := s.ctrl.SetWorkers(id, workers); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetImage(w http.ResponseWriter, id string) {
	// TODO
	img, err := s.ctrl.GetImage(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := png.Encode(w, img); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
