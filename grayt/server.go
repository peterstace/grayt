package grayt

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"sync"
	"sync/atomic"

	uuid "github.com/satori/go.uuid"
)

type scene struct {
	sceneFn     func() Scene
	ascpectWide int
	ascpectHigh int
}

type Server struct {
	scenes    map[string]scene
	resources []*resource
}

func NewServer() *Server {
	return &Server{
		scenes:    make(map[string]scene),
		resources: []*resource{}, // Get's serialised to JSON, so important it's not nil.
	}
}

func (s *Server) ListenAndServe(addr string) error {
	http.Handle("/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/scenes", s.handleGetScenesCollection)
	http.HandleFunc("/renders", s.handleRendersCollection)

	log.Printf("Listening for HTTP on %v", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) Register(
	name string,
	SkyFn func(Vector) Colour,
	cam CameraBlueprint,
	objFn func() ObjectList,
) {
	s.scenes[name] = scene{
		sceneFn: func() Scene {
			return Scene{
				Camera:  cam,
				Objects: objFn(),
				Sky:     SkyFn,
			}
		},
		ascpectWide: cam.aspectWide,
		ascpectHigh: cam.aspectHigh,
	}
}

type resource struct {
	sync.Mutex
	uuid   uuid.UUID
	render *render
	cancel func() // set to nil if the render isn't running

	scene     Scene
	sceneName string

	pxWide, pxHigh int
}

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func internalError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) handleGetScenesCollection(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Code       string `json:"code"`
		AspectWide int    `json:"aspect_wide"`
		AspectHigh int    `json:"aspect_high"`
	}
	var responses []response
	for name, scene := range s.scenes {
		responses = append(responses, response{name, scene.ascpectWide, scene.ascpectHigh})
	}
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		internalError(w, err)
	}
}

func (s *Server) handleRendersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var form struct {
			Scene  string `json:"scene"`
			PxWide int    `json:"px_wide"`
		}
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, fmt.Sprintf("could not decode form: %v", err), http.StatusBadRequest)
			return
		}
		sceneInfo, ok := s.scenes[form.Scene]
		if !ok {
			http.Error(w, fmt.Sprintf("scene %q not found", form.Scene), http.StatusBadRequest)
			return
		}
		if form.PxWide == 0 {
			http.Error(w, "px_wide not set", http.StatusBadRequest)
			return
		}

		id := uuid.Must(uuid.NewV4())
		rsrc := &resource{
			uuid:      id,
			scene:     sceneInfo.sceneFn(), // TODO: This could take some time.
			sceneName: form.Scene,
			pxWide:    form.PxWide,
			pxHigh:    form.PxWide * sceneInfo.ascpectHigh / sceneInfo.ascpectWide,
		}
		s.resources = append(s.resources, rsrc)

		fmt.Fprintf(w, `{"uuid":%q}`, id)
		http.HandleFunc("/renders/"+id.String()+"/image", rsrc.handleGetImage)
		http.HandleFunc("/renders/"+id.String()+"/running", rsrc.handlePutRunning)
	case http.MethodGet:
		type props struct {
			ID        uuid.UUID `json:"uuid"`
			Running   bool      `json:"running"`
			Scene     string    `json:"scene"`
			Completed uint64    `json:"completed"`
			Passes    uint64    `json:"passes"`
			PxWide    int       `json:"px_wide"`
			PxHigh    int       `json:"px_high"`
		}
		propList := []props{} // Populate as empty array since it goes to JSON.
		for _, rsrc := range s.resources {
			rsrc.Lock()
			defer rsrc.Unlock()

			var completed, passes uint64
			if rsrc.render != nil {
				completed = atomic.LoadUint64(&rsrc.render.completed)
				passes = atomic.LoadUint64(&rsrc.render.passes)
			}

			propList = append(propList, props{
				rsrc.uuid,
				rsrc.render != nil,
				rsrc.sceneName,
				completed,
				passes,
				rsrc.pxWide,
				rsrc.pxHigh,
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
	rsrc.Lock()
	defer rsrc.Unlock()

	if rsrc.render == nil {
		img := image.NewGray(image.Rect(0, 0, 320, 240))
		if err := png.Encode(w, img); err != nil {
			internalError(w, err)
		}
		return
	}

	// Disable caching, since this image will update often.
	w.Header().Set("Cache-Control", "no-cache")

	img := rsrc.render.accum.toImage(1.0)
	if err := png.Encode(w, img); err != nil {
		internalError(w, err)
		return
	}
}

func (rsrc *resource) handlePutRunning(w http.ResponseWriter, r *http.Request) {
	rsrc.Lock()
	defer rsrc.Unlock()

	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		internalError(w, err)
		return
	}

	b, err := strconv.ParseBool(string(buf))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if (rsrc.cancel != nil) == b {
		// Already in the correct state.
		return
	}

	if !b {
		rsrc.cancel()
		rsrc.cancel = nil
		return
	}

	if rsrc.render == nil {
		pxHigh := rsrc.pxWide * rsrc.scene.Camera.aspectHigh / rsrc.scene.Camera.aspectWide
		rsrc.render = newRender(
			rsrc.pxWide,
			rsrc.scene,
			newAccumulator(rsrc.pxWide, pxHigh),
		)
	}
	var ctx context.Context
	ctx, rsrc.cancel = context.WithCancel(context.Background())
	go func() {
		rsrc.render.traceImage(ctx)
	}()
}
