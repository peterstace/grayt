package grayt

import (
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"path/filepath"
	"strconv"

	gscene "github.com/peterstace/grayt/scene"
	"github.com/peterstace/grayt/trace"
	uuid "github.com/satori/go.uuid" // TODO: don't use uuids...
)

type scene struct {
	sceneFn func() trace.Scene
}

type Server struct {
	resources []*resource
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Load(storageDir string) error {
	entries, err := ioutil.ReadDir(storageDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fname := filepath.Join(storageDir, entry.Name())
		id, err := uuid.FromString(entry.Name())
		if err != nil {
			continue // it's just some other file
		}
		sceneName, acc, err := loadAccumulator(fname)
		if err != nil {
			return fmt.Errorf("could not load from file %q: %v", fname, err)
		}
		scene, err := s.lookupScene(sceneName)
		if err != nil {
			return err
		}
		s.addResource(id, sceneName, scene, acc)
	}
	return nil
}

func (s *Server) lookupScene(name string) (trace.Scene, error) {
	// TODO: allow address to be configured
	resp, err := http.Get("http://scenelib:80/scene?name=" + url.QueryEscape(name))
	if err != nil {
		return trace.Scene{}, fmt.Errorf("fetching scene: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return trace.Scene{}, fmt.Errorf("fetching scene: %s", resp.Status)
	}
	var scn gscene.Scene
	if err := json.NewDecoder(resp.Body).Decode(&scn); err != nil {
		return trace.Scene{}, fmt.Errorf("decoding scene: %v", err)
	}
	return trace.BuildScene(scn), nil
}

func (s *Server) Save(storageDir string) error {
	for _, rsrc := range s.resources {
		fname := filepath.Join(storageDir, rsrc.uuid.String())
		if err := rsrc.render.accum.save(fname, rsrc.sceneName); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) ListenAndServe(addr string) error {
	http.Handle("/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/scenes", s.handleGetScenesCollection)
	http.HandleFunc("/renders", s.handleRendersCollection)

	log.Printf("listening for http on %v", addr)
	return http.ListenAndServe(addr, nil)
}

type resource struct {
	uuid      uuid.UUID
	render    *render
	sceneName string
}

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func internalError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) handleGetScenesCollection(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Code string `json:"code"`
	}
	// TODO: should get the full scene list from the grayt.scenelib server.
	responses := []response{
		{Code: "cornellbox_classic"},
		{Code: "cornellbox_splitbox"},
		{Code: "cornellbox_mirror"},
		{Code: "cornellbox_spheretree"},
	}
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		internalError(w, err)
	}
}

func (s *Server) addResource(id uuid.UUID, sceneName string, scene trace.Scene, acc *accumulator) {
	rsrc := &resource{
		uuid:      id,
		sceneName: sceneName,
	}
	rsrc.render = newRender(
		scene,
		acc,
	)
	go rsrc.render.traceImage()
	http.HandleFunc("/renders/"+id.String()+"/image", rsrc.handleGetImage)
	http.HandleFunc("/renders/"+id.String()+"/workers", rsrc.handlePutWorkers)
	s.resources = append(s.resources, rsrc)
}

// TODO: Break this method into two parts, for GET and POST.
func (s *Server) handleRendersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var form struct {
			Scene  string `json:"scene"`
			PxWide int    `json:"px_wide"`
			PxHigh int    `json:"px_high"`
		}
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, fmt.Sprintf("could not decode form: %v", err), http.StatusBadRequest)
			return
		}
		if form.PxWide == 0 || form.PxHigh == 0 {
			http.Error(w, "px_wide or px_high not set", http.StatusBadRequest)
			return
		}
		scene, err := s.lookupScene(form.Scene)
		if err != nil {
			http.Error(w, "scene lookup failed: "+err.Error(), http.StatusBadRequest)
			return
		}
		id := uuid.Must(uuid.NewV4())
		s.addResource(id, form.Scene, scene, newAccumulator(form.PxWide, form.PxHigh))
		fmt.Fprintf(w, `{"uuid":%q}`, id)

	case http.MethodGet:
		type props struct {
			ID               uuid.UUID `json:"uuid"`
			Scene            string    `json:"scene"`
			Completed        string    `json:"completed"`
			Passes           int64     `json:"passes"`
			PxWide           int       `json:"px_wide"`
			PxHigh           int       `json:"px_high"`
			RequestedWorkers int64     `json:"requested_workers"`
			ActualWorkers    int64     `json:"actual_workers"`
			TraceRate        string    `json:"trace_rate"`
		}
		propList := []props{} // Populate as empty array since it goes to JSON.
		for _, rsrc := range s.resources {
			status := rsrc.render.status()
			propList = append(propList, props{
				rsrc.uuid,
				rsrc.sceneName,
				displayFloat64(float64(status.completed)),
				status.passes,
				rsrc.render.accum.wide,
				rsrc.render.accum.high,
				status.requestedWorkers,
				status.actualWorkers,
				displayFloat64(float64(status.traceRate)) + " Hz",
			})
		}
		if err := json.NewEncoder(w).Encode(propList); err != nil {
			internalError(w, err)
		}
	default:
		writeError(w, http.StatusMethodNotAllowed)
	}
}

func (rsrc *resource) handleGetImage(w http.ResponseWriter, r *http.Request) {
	img := rsrc.render.accum.toImage(1.0)
	if err := png.Encode(w, img); err != nil {
		internalError(w, err)
		return
	}
}

func (rsrc *resource) handlePutWorkers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		internalError(w, err)
		return
	}
	workers, err := strconv.Atoi(string(buf))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if workers < 0 {
		return // Ignore requests for negative workers.
	}

	rsrc.render.setWorkers(int64(workers))
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
