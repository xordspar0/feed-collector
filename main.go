package main

import (
	"neolog.xyz/feed-collector/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	s := server.New()
	err := s.Start()
	if err != nil {
		log.Error(err.Error())
	}
}
