package grayt

import (
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"sync"
)

func ListenAndServe(addr string) error {
	http.Handle("/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/scenes", handleGetScenesCollection)
	http.HandleFunc("/renders", handlePostRendersCollection)

	log.Printf("Listening for HTTP on %v", addr)
	return http.ListenAndServe(addr, nil)
}

/*
	POST  /renders                    - Adds a new render resource, not started, with default settings.
	GET   /renders/{uuid}             - Gets all information about the render resource.
	PUT   /renders/{uuid}/{property}  - Sets property of the render.
	GET   /renders/{uuid}/image       - Creates an image.
*/

// TODO: Rename?
type resource struct {
	sync.Mutex
	uuid string
	render
	cancel func() // set to nil if the render isn't running
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

var lastUUID int

func handlePostRendersCollection(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %q\n", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	//uuid := fmt.Sprintf("%d", time.Now().UnixNano())
	lastUUID++
	uuid := fmt.Sprintf("%d", lastUUID)

	rsrc := &resource{uuid: uuid}

	fmt.Fprintf(w, `{"uuid":%q}`, uuid)

	http.HandleFunc("/renders/"+uuid, rsrc.handleGetAll)
	http.HandleFunc("/renders/"+uuid+"/image", rsrc.handleGetImage)
	http.HandleFunc("/renders/"+uuid+"/scene", rsrc.handlePutScene)
	http.HandleFunc("/renders/"+uuid+"/running", rsrc.handlePutRunning)
}

func (rsrc *resource) handleGetAll(w http.ResponseWriter, r *http.Request) {
	rsrc.Lock()
	defer rsrc.Unlock()

	fmt.Fprintf(w, `{"uuid":%q}`, rsrc.uuid)
	// TODO: Add other properties to the response.
}

func (rsrc *resource) handleGetImage(w http.ResponseWriter, r *http.Request) {
	rsrc.Lock()
	defer rsrc.Unlock()

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

	sceneFn, ok := scenes[string(buf)]
	if !ok {
		http.Error(w, fmt.Sprintf("scene %q not found", string(buf)), http.StatusBadRequest)
		return
	}

	rsrc.scene = sceneFn() // TODO: This function could take some time...
	fmt.Printf("%p\n", &rsrc.scene.Camera)
	fmt.Printf("%+v\n", rsrc.scene.Camera)
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
	} else {
		const pxWide = 320
		rsrc.render.pxWide = pxWide
		rsrc.render.numWorkers = 1
		high := pxWide * rsrc.scene.Camera.aspectHigh / rsrc.scene.Camera.aspectWide
		rsrc.render.accum = newAccumulator(pxWide, high)

		var ctx context.Context
		ctx, rsrc.cancel = context.WithCancel(context.Background())
		go func() {
			rsrc.render.traceImage(ctx)
		}()
	}
}
