package api

import (
	"log"
	"net/http"
	"time"
)

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf(`--> %s %s`, req.Method, req.URL)
		start := time.Now()
		wrappedW := wrappedResponseWriter{ResponseWriter: w}
		next.ServeHTTP(&wrappedW, req)
		elapsed := time.Since(start)
		log.Printf(`<-- %d %s`, wrappedW.statusCode, elapsed)
	})
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	statusCode  int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
	w.wroteHeader = true
}

func (w *wrappedResponseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}
