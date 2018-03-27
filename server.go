package main

import (
	"neolog.xyz/feed-collector/feeds"
	"neolog.xyz/feed-collector/nextcloudnews"
	"neolog.xyz/feed-collector/pocket"

	"github.com/rs/cors"
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

	NextcloudNewsHost     string
	NextcloudNewsUser     string
	NextcloudNewsPassword string
	PocketAccessToken     string
	PocketConsumerKey     string
}

var healthyResponse string = fmt.Sprintf(`{"status": "%d"}`, http.StatusOK)

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc(path.Join(s.RootEndpoint, "health"), s.health)
	mux.HandleFunc(path.Join(s.RootEndpoint, "feeds"), s.feeds)

	corsMux := cors.Default().Handler(mux)
	return http.ListenAndServe(":"+s.Port, corsMux)
}

func (s *Server) health(w http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{
		"host": req.Host,
		"url":  req.URL,
	})
	logger.Debug("Handling request")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, healthyResponse)

	logger.WithFields(log.Fields{
		"status": http.StatusOK,
	}).Info("Received request")
}

func (s *Server) feeds(w http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{
		"host": req.Host,
		"url":  req.URL,
	})
	logger.Debug("Handling request")

	responseData := feeds.New()

	var feed *feeds.Feed
	var err error

	feed = responseData.AddFeed("Nextcloud News", s.NextcloudNewsHost + "/apps/news/")
	if s.NextcloudNewsUser != "" && s.NextcloudNewsPassword != "" {
		logger.WithFields(log.Fields{
			"feed": feed.Name,
		}).Debug("Requesting feed unread count")

		nextcloudnewsHost := nextcloudnews.New(
			feed.URL,
			s.NextcloudNewsUser,
			s.NextcloudNewsPassword,
		)

		feed.Count, err = nextcloudnewsHost.GetUnreadCount()
		if err != nil {
			logger.WithFields(log.Fields{
				"feed":  feed.Name,
				"error": err.Error(),
			}).Error("Failed to count articles")
		}

		logger.WithFields(log.Fields{
			"feed":  feed.Name,
			"count": feed.Count,
		}).Debug("Received unread count")
	} else {
		logger.WithFields(
			log.Fields{"feed": feed.Name},
		).Error("Missing credentials for feed")
	}

	feed = responseData.AddFeed("Pocket", "https://getpocket.com")
	if s.PocketAccessToken == "" {
		logger.WithFields(
			log.Fields{"feed": feed.Name},
		).Error("Missing Pocket access token")
	} else if s.PocketConsumerKey == "" {
		logger.WithFields(
			log.Fields{"feed": feed.Name},
		).Error("Missing Pocket consumer key")
	} else {
		logger.WithFields(log.Fields{
			"feed": feed.Name,
		}).Debug("Requesting feed unread count")

		pocketUser := pocket.NewUser(s.PocketAccessToken, s.PocketConsumerKey)
		feed.Count, err = pocketUser.GetUnreadCount()
		if err != nil {
			logger.WithFields(log.Fields{
				"feed":  feed.Name,
				"error": err.Error(),
			}).Error("Failed to count articles")
		}
	}

	resp, err := json.Marshal(responseData)
	if err != nil {
		logger.Error("Failed to marshal json response: " + err.Error())
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(resp)

	logger.WithFields(log.Fields{
		"status": http.StatusOK,
	}).Info("Served request")
}
