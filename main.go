package main

import (
	"neolog.xyz/feed-collector/server"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"fmt"
	"os"
)

var servername = "feed-collector"
var version = "devel"

func main() {
	app := cli.NewApp()
	app.Name = servername
	app.Usage = "An API that collects your news feeds"
	app.Action = run
	app.Version = version
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port, p",
			Value:  80,
			Usage:  "The port to run the server on",
			EnvVar: "SQUIRRELBOT_PORT",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}

func run(c *cli.Context) {
	log.Info(fmt.Sprintf("Starting %s version %s", servername, version))
	s := server.New(c.String("port"), "/")
	err := s.Start()
	if err != nil {
		log.Error(err.Error())
	}
}
