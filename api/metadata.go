package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peterstace/grayt/xmath"
)

type metadata struct {
	SceneName string           `json:"scene_name"`
	Created   time.Time        `json:"created"`
	Dim       xmath.Dimensions `json:"dim"`
}

func saveMetadata(m metadata, filename string) error {
	buf, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("could not save metadata: %v", err)
	}
	return ioutil.WriteFile(filename, buf, 0664)
}

func (s *Server) loadRenders() error {
	fileInfos, err := ioutil.ReadDir(s.dataDir)
	if err != nil {
		return fmt.Errorf("could not read data dir: %v", err)
	}
	for _, fi := range fileInfos {
		fname := filepath.Join(s.dataDir, fi.Name())
		if fi.IsDir() || filepath.Ext(fname) != ".json" {
			continue
		}
		var m metadata
		f, err := os.Open(fname)
		if err != nil {
			return fmt.Errorf("could not open metadata file: %v", err)
		}
		if err := json.NewDecoder(f).Decode(&m); err != nil {
			return fmt.Errorf("could not load metadata: %v", err)
		}
		f.Close()

		id := strings.TrimSuffix(filepath.Base(fname), ".json")
		accumFilename := filepath.Join(filepath.Dir(fname), id+".data")
		if err := s.ctrl.newRender(
			id, accumFilename, m.Created, m.SceneName, m.Dim,
		); err != nil {
			return fmt.Errorf("could not create render: %v", err)
		}
	}
	return nil
}
