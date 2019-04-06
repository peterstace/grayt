package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
