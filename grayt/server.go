package grayt

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync/atomic"
)

func ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	s := new(server)
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/status", s.handleStatus)

	// Run scene in background. TODO Eventually this will be trigged using
	// a /start endpoint.
	go func() {
		if err := RunScene(&s.completed); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	return http.ListenAndServe(addr, mux)
}

type server struct {
	completed uint64
}

func (s *server) handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func (s *server) handleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"completed":%d}`, atomic.LoadUint64(&s.completed))
}
