package main

import (
	"neolog.xyz/feed-collector/feeds"
	"neolog.xyz/feed-collector/nextcloudnews"
	"neolog.xyz/feed-collector/pocket"

	"github.com/rs/cors"

	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"
)

type Server struct {
	Port         string
	RootEndpoint string
	Logger       *slog.Logger

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
	logger := s.Logger.With(
		"host", req.Host,
		"url", req.URL,
	)
	logger.Debug("Handling request")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, healthyResponse)

	logger.Debug("Received request", "status", http.StatusOK)
}

func (s *Server) feeds(w http.ResponseWriter, req *http.Request) {
	logger := s.Logger.With(
		"host", req.Host,
		"url", req.URL,
	)
	logger.Debug("Handling request")

	responseData := feeds.New()

	var feed *feeds.Feed
	var err error

	feed = responseData.AddFeed("Nextcloud News", s.NextcloudNewsHost+"/apps/news/")
	if s.NextcloudNewsUser != "" && s.NextcloudNewsPassword != "" {
		logger.Debug("Requesting feed unread count",
			"feed", feed.Name,
		)

		nextcloudnewsHost := nextcloudnews.New(
			s.NextcloudNewsHost,
			s.NextcloudNewsUser,
			s.NextcloudNewsPassword,
		)

		feed.Count, err = nextcloudnewsHost.GetUnreadCount()
		if err != nil {
			logger.Error("Failed to count articles",
				"feed", feed.Name,
				"error", err.Error(),
			)
		}

		logger.Debug("Received unread count",
			"feed", feed.Name,
			"count", feed.Count,
		)
	} else {
		logger.Error("Missing credentials for feed", "feed", feed.Name)
	}

	feed = responseData.AddFeed("Pocket", "https://getpocket.com")
	if s.PocketAccessToken == "" {
		logger.Error("Missing Pocket access token", "feed", feed.Name)
	} else if s.PocketConsumerKey == "" {
		logger.Error("Missing Pocket consumer key", "feed", feed.Name)
	} else {
		logger.Debug("Requesting feed unread count", "feed", feed.Name)

		pocketUser := pocket.NewUser(s.PocketAccessToken, s.PocketConsumerKey)
		feed.Count, err = pocketUser.GetUnreadCount()
		if err != nil {
			logger.Error("Failed to count articles",
				"feed", feed.Name,
				"error", err.Error(),
			)
		}
	}

	resp, err := json.Marshal(responseData)
	if err != nil {
		logger.Error("Failed to marshal json response", "error", err.Error())
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(resp)

	logger.Info("Served request",
		"status", http.StatusOK,
	)
}
