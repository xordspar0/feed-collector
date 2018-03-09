package server

import (
	"neolog.xyz/feed-collector/feeds"

	log "github.com/sirupsen/logrus"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

type Server struct {
	Port         string
	RootEndpoint string
}

var healthyResponse string = fmt.Sprintf(`{"status": "%d"}`, http.StatusOK)

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
	io.WriteString(w, healthyResponse)
}

func (s *Server) feeds(w http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{
		"host": req.Host,
		"url":  req.URL,
	})

	feeds, err := json.Marshal(feeds.New())
	if err != nil {
		logger.Error("Failed to marshal json response: " + err.Error())
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}

	logger.Info()
	w.Write(feeds)
}
