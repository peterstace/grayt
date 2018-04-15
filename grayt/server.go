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
	"time"

	uuid "github.com/satori/go.uuid"
)

type scene struct {
	sceneFn     func() Scene
	ascpectWide int
	ascpectHigh int
}

type Server struct {
	scenes map[string]scene
	idList []uuid.UUID
}

func NewServer() *Server {
	return &Server{
		scenes: make(map[string]scene),
		idList: []uuid.UUID{}, // Get's serialised to JSON, so important it's not nil.
	}
}

func (s *Server) ListenAndServe(addr string) error {
	http.Handle("/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/scenes", middleware(s.handleGetScenesCollection))
	http.HandleFunc("/renders", middleware(s.handleRendersCollection))

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

func middleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(`info="handling request" method=%q path=%q`, r.Method, r.URL.Path)
		fn(w, r)
		log.Printf(`info="finished" duration=%q`, time.Since(start))
	}
}

/*
	GET   /scenes                     - Gets a list of all available scenes.
	POST  /renders                    - Adds a new render resource, not started, with default settings.
	GET   /renders                    - Gets a list of all existing render UUIDs.
	GET   /renders/{uuid}             - Gets all information about the render resource.
	PUT   /renders/{uuid}/scene       - Sets the scene property of the render resource.
	PUT   /renders/{uuid}/running     - Sets the runnig property of the render resource.
	GET   /renders/{uuid}/image       - Creates an image.
*/

type resource struct {
	sync.Mutex
	uuid   uuid.UUID
	render *render
	cancel func() // set to nil if the render isn't running

	scene     Scene
	sceneName string

	s *Server // TODO: This feels like a hack.
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
			Scene string `json:"scene"`
			// TODO: screen resolutions go here as well.
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

		id := uuid.Must(uuid.NewV4())
		s.idList = append(s.idList, id)
		rsrc := &resource{
			uuid:      id,
			scene:     sceneInfo.sceneFn(), // TODO: This could take some time.
			sceneName: form.Scene,
			s:         s,
		}
		fmt.Fprintf(w, `{"uuid":%q}`, id)
		http.HandleFunc("/renders/"+id.String(), middleware(rsrc.handleGetAll))
		http.HandleFunc("/renders/"+id.String()+"/image", middleware(rsrc.handleGetImage))
		http.HandleFunc("/renders/"+id.String()+"/running", middleware(rsrc.handlePutRunning))
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(s.idList); err != nil {
			internalError(w, err)
		}
	default:
		writeError(w, http.StatusMethodNotAllowed)
	}
}

func (rsrc *resource) handleGetAll(w http.ResponseWriter, r *http.Request) {
	rsrc.Lock()
	defer rsrc.Unlock()

	var completed, passes uint64
	if rsrc.render != nil {
		completed = atomic.LoadUint64(&rsrc.render.completed)
		passes = atomic.LoadUint64(&rsrc.render.passes)
	}

	props := struct {
		ID        uuid.UUID `json:"uuid"`
		Running   bool      `json:"running"`
		Scene     string    `json:"scene"`
		Completed uint64    `json:"completed"`
		Passes    uint64    `json:"passes"`
	}{
		rsrc.uuid,
		rsrc.render != nil,
		rsrc.sceneName,
		completed,
		passes,
	}

	if err := json.NewEncoder(w).Encode(props); err != nil {
		internalError(w, err)
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
		// TODO: sceneFunc might not exist yet.
		const pxWide = 320
		pxHigh := pxWide * rsrc.scene.Camera.aspectHigh / rsrc.scene.Camera.aspectWide
		rsrc.render = &render{
			pxWide:     pxWide,
			numWorkers: 1,
			scene:      rsrc.scene,
			accum:      newAccumulator(pxWide, pxHigh),
		}
	}
	var ctx context.Context
	ctx, rsrc.cancel = context.WithCancel(context.Background())
	go func() {
		rsrc.render.traceImage(ctx)
	}()
}
