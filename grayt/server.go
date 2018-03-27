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

func ListenAndServe(addr string) error {
	http.Handle("/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/scenes", middleware(handleGetScenesCollection))
	http.HandleFunc("/renders", middleware(handleRendersCollection))

	log.Printf("Listening for HTTP on %v", addr)
	return http.ListenAndServe(addr, nil)
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

	sceneFunc func() Scene
	sceneName string
}

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func internalError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func handleGetScenesCollection(w http.ResponseWriter, r *http.Request) {
	var ss []string
	for s := range scenes {
		ss = append(ss, s)
	}
	if err := json.NewEncoder(w).Encode(ss); err != nil {
		internalError(w, err)
	}
}

var idList []uuid.UUID = []uuid.UUID{} // TODO: Don't be a global.

func handleRendersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		id := uuid.Must(uuid.NewV4())
		idList = append(idList, id)
		rsrc := &resource{uuid: id}
		fmt.Fprintf(w, `{"uuid":%q}`, id)
		http.HandleFunc("/renders/"+id.String(), middleware(rsrc.handleGetAll))
		http.HandleFunc("/renders/"+id.String()+"/image", middleware(rsrc.handleGetImage))
		http.HandleFunc("/renders/"+id.String()+"/scene", middleware(rsrc.handlePutScene))
		http.HandleFunc("/renders/"+id.String()+"/running", middleware(rsrc.handlePutRunning))
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(idList); err != nil {
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

	img := rsrc.render.accum.toImage(1.0)
	if err := png.Encode(w, img); err != nil {
		internalError(w, err)
		return
	}

}

func (rsrc *resource) handlePutScene(w http.ResponseWriter, r *http.Request) {
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

	sceneName := string(buf)
	sceneFunc, ok := scenes[sceneName]
	if !ok {
		http.Error(w, fmt.Sprintf("scene %q not found", string(buf)), http.StatusBadRequest)
		return
	}
	rsrc.sceneFunc = sceneFunc
	rsrc.sceneName = sceneName
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
		scene := rsrc.sceneFunc() // TODO: This could take some time to run.
		const pxWide = 320
		pxHigh := pxWide * scene.Camera.aspectHigh / scene.Camera.aspectWide
		rsrc.render = &render{
			pxWide:     pxWide,
			numWorkers: 1,
			scene:      scene,
			accum:      newAccumulator(pxWide, pxHigh),
		}
	}
	var ctx context.Context
	ctx, rsrc.cancel = context.WithCancel(context.Background())
	go func() {
		rsrc.render.traceImage(ctx)
	}()
}
