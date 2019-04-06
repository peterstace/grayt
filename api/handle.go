package api

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"image/png"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"time"

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
	renders := s.ctrl.getRenders()
	sort.Slice(renders, func(i, j int) bool {
		return renders[i].Created.Before(renders[j].Created)
	})
	if err := json.NewEncoder(w).Encode(renders); err != nil {
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

	id := generateID()
	accumFilename := filepath.Join(s.dataDir, id+".data")
	now := time.Now()
	dim := xmath.Dimensions{form.PxWide, form.PxHigh}
	metadataFilename := filepath.Join(s.dataDir, id+".json")
	if err := saveMetadata(metadata{form.Scene, now, dim}, metadataFilename); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.ctrl.newRender(id, accumFilename, now, form.Scene, dim); err != nil {
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

	if err := s.ctrl.setWorkers(id, workers); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func generateID() string {
	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(time.Now().Unix()))
	sum := crc64.Checksum(buf[:], crc64.MakeTable(crc64.ECMA))
	return fmt.Sprintf("%X", sum)
}

func (s *Server) handleGetImage(w http.ResponseWriter, id string) {
	img, err := s.ctrl.getImage(id)
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
