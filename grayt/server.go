package grayt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"regexp"
	"strconv"
	"time"
)

func ListenAndServe(addr string) error {
	s := &server{
		renders: make(map[string]*render),
	}

	http.Handle("/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/renders", s.handleEndpoint)
	http.HandleFunc("/renders/", s.handleEndpoint)

	log.Printf("Listening for HTTP on %v", addr)
	return http.ListenAndServe(addr, nil)
}

type server struct {
	renders map[string]*render
}

var endpointRE = regexp.MustCompile("^/renders/(?:([0-9]+)(?:/([a-z]+))?)?$")

func (s *server) handleEndpoint(w http.ResponseWriter, r *http.Request) {
	/*
		POST  /renders                    - Adds a new render resource, not started, with default settings.
		GET   /renders/{uuid}             - Gets all information about the render resource.
		PUT   /renders/{uuid}/{property}  - Sets property of the render.
		GET   /renders/{uuid}/image       - Creates an image.
	*/

	log.Printf("%s %q\n", r.Method, r.URL.Path)

	captured := endpointRE.FindStringSubmatch(r.URL.Path)
	if len(captured) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if len(captured) != 3 {
		// We don't expect to see this. Should either capture 0 or 3 strings from the regexp.
		fmt.Fprintf(w, "Something broke internally while parsing url %q, please file a bug report.", r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uuid := captured[1]
	prop := captured[2]
	if uuid == "" {
		// "/renders"
		if r.Method == http.MethodPost {
			s.handlePostRenders(w)
			return
		}
	} else if prop == "" {
		// "/renders/{uuid}
		if r.Method == http.MethodGet {
			s.handleGetRender(w, uuid)
			return
		}
	} else {
		// "/renders/{uuid}/{property}
		if r.Method == http.MethodPut {
			val, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "Could not ready body: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			s.handlePutRender(w, uuid, prop, string(val))
			return
		} else if r.Method == http.MethodGet && prop == "image" {
			s.handleGetRenderImage(w, uuid)
			return
		}
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *server) handlePostRenders(w http.ResponseWriter) {
	r := &render{
		completed:  0,
		PxWide:     320,
		NumWorkers: 1,
		Scene:      "", // Must be set, no default.
	}
	uuid := fmt.Sprintf("%d", time.Now().UnixNano())
	s.renders[uuid] = r
	fmt.Fprintf(w, `{"uuid":%q}`, uuid)
}

func (s *server) handleGetRender(w http.ResponseWriter, uuid string) {
	r, ok := s.renders[uuid]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	type response struct {
		completed int `json:"completed"`
	}
	buf, err := json.Marshal(struct {
		uuid      string `json:"uuid"`
		completed int    `json:"completed"`
	}{
		uuid,
		int(r.completed),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(buf)
}

func toInt(w http.ResponseWriter, val string) (int, bool) {
	v, err := strconv.Atoi(val)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return 0, false
	}
	return v, true
}

func (s *server) handlePutRender(w http.ResponseWriter, uuid, prop, propVal string) {
	r, ok := s.renders[uuid]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch prop {
	case "px_wide":
		v, ok := toInt(w, propVal)
		if !ok {
			return
		}
		r.PxWide = v
	case "quality":
		v, ok := toInt(w, propVal)
		if !ok {
			return
		}
		r.Quality = v
	case "num_workers":
		v, ok := toInt(w, propVal)
		if !ok {
			return
		}
		r.NumWorkers = v
	case "scene":
		_, ok := scenes[propVal]
		if !ok {
			fmt.Fprintf(w, "Unknown scene: %v", propVal)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Scene = propVal
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *server) handleGetRenderImage(w http.ResponseWriter, uuid string) {
	// TODO
}

func checkMethod(w http.ResponseWriter, r *http.Request, allowedMethods ...string) bool {
	for _, allowed := range allowedMethods {
		if r.Method == allowed {
			return true
		}
	}
	fmt.Fprintf(w, "only allowed: %v", allowedMethods)
	w.WriteHeader(http.StatusMethodNotAllowed)
	return false
}
