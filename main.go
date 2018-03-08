package main

import (
	"neolog.xyz/feed-collector/server"

	log "github.com/sirupsen/logrus"

	"fmt"
)

var servername = "feed-collector"
var version = "devel"

func main() {
	log.Info(fmt.Sprintf("Starting %s version %s", servername, version))
	s := server.New()
	err := s.Start()
	if err != nil {
		log.Error(err.Error())
	}
}
