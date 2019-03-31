package api

import "net/http"

func NewServer() *Server {
	return &Server{}
}

type Server struct {
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	panic("not implemented")
}
