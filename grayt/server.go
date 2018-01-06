package grayt

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
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
	uuid string
	render
	img image.Image
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

func handlePostRendersCollection(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %q\n", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	uuid := fmt.Sprintf("%d", time.Now().UnixNano())
	rsrc := &resource{
		uuid: uuid,
		render: render{
			completed:  0,
			pxWide:     320,
			quality:    10,
			numWorkers: 1,
		},
	}
	fmt.Fprintf(w, `{"uuid":%q}`, uuid)

	http.HandleFunc("/renders/"+uuid, rsrc.handleGetAll)
	http.HandleFunc("/renders/"+uuid+"/image", rsrc.handleGetImage)
	http.HandleFunc("/renders/"+uuid+"/scene", rsrc.handlePutScene)
	http.HandleFunc("/renders/"+uuid+"/running", rsrc.handlePutRunning)
}

func (rsrc *resource) handleGetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"uuid":%q}`, rsrc.uuid)
	// TODO: Add other properties to the response.
}

func (rsrc *resource) handleGetImage(w http.ResponseWriter, r *http.Request) {
	if err := png.Encode(w, rsrc.img); err != nil {
		internalError(w, err)
	}
}

func (rsrc *resource) handlePutScene(w http.ResponseWriter, r *http.Request) {
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
}

func (rsrc *resource) handlePutRunning(w http.ResponseWriter, r *http.Request) {
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

	// TODO: Handle starting and stopping properly.
	if b {
		go func() {
			acc := new(accumulator)
			pxHigh := rsrc.pxWide * rsrc.scene.Camera.aspectHigh / rsrc.scene.Camera.aspectWide
			n := rsrc.pxWide * pxHigh
			acc.pixels = make([]Colour, n)
			acc.wide = rsrc.pxWide
			acc.high = pxHigh

			rsrc.img = rsrc.render.traceImage(acc)
		}()
	}
}
