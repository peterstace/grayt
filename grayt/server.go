package grayt

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

func ListenAndServe(addr string) error {
	http.Handle("/", http.FileServer(http.Dir("assets")))
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

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func handlePostRendersCollection(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %q\n", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	uuid := fmt.Sprintf("%d", time.Now().UnixNano())
	rn := &render{
		completed:  0,
		PxWide:     320,
		Quality:    10,
		NumWorkers: 0,  // Rendering starts when this is set greater than 0.
		Scene:      "", // Must be set, no default.
		UUID:       uuid,
	}
	fmt.Fprintf(w, `{"uuid":%q}`, uuid)

	http.HandleFunc("/renders/"+uuid, rn.handleGetAll)
	http.HandleFunc("/renders/"+uuid+"/image", rn.handleGetImage)
	http.HandleFunc("/renders/"+uuid+"/px_wide", rn.handlePutPxWide)
	http.HandleFunc("/renders/"+uuid+"/quality", rn.handlePutQuality)
	http.HandleFunc("/renders/"+uuid+"/num_workers", rn.handlePutNumWorkers)
	http.HandleFunc("/renders/"+uuid+"/scene", rn.handlePutScene)
}

func (rn *render) handleGetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"uuid":%q}`, rn.UUID)
	// TODO: Add other properties to the response.
}

func (rn *render) handleGetImage(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Will need to stop the tracer, then get the image?
}

func (rn *render) handlePutPxWide(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (rn *render) handlePutQuality(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	q, err := strconv.Atoi(string(buf))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Some kind of locking around this?
	rn.Quality = q
}

func (rn *render) handlePutNumWorkers(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (rn *render) handlePutScene(w http.ResponseWriter, r *http.Request) {
	// TODO
}
