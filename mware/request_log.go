package mware

import (
	"log"
	"net/http"
	"time"
)

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// TODO: Request ID
		log.Printf(`Info="HTTP Request" Method=%q URL=%q`, req.Method, req.URL)
		start := time.Now()
		next.ServeHTTP(w, req)
		elapsed := time.Since(start)
		log.Printf(`Info="HTTP Response" Duration=%q`, elapsed)
	})
}
