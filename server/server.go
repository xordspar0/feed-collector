package server

import (
	log "github.com/sirupsen/logrus"

	"fmt"
	"io"
	"net/http"
	"path"
)

type Server struct {
	Port         string
	RootEndpoint string
}

func New(port string, rootEndpoint string) Server {
	return Server{
		Port:         port,
		RootEndpoint: rootEndpoint,
	}
}

func (s *Server) Start() error {
	http.HandleFunc(path.Join(s.RootEndpoint, "health"), s.health)
	http.HandleFunc(path.Join(s.RootEndpoint, "feeds"), s.feeds)

	return http.ListenAndServe(":"+s.Port, nil)
}

func (s *Server) health(w http.ResponseWriter, req *http.Request) {
	log.WithFields(log.Fields{
		"host": req.Host,
		"url":  req.URL,
	}).Info()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, fmt.Sprintf(`{"status": "%d"}`, http.StatusOK))
}

func (s *Server) feeds(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, `{"feeds": []}`)
}
