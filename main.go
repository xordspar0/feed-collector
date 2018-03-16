package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

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
			EnvVar: "FEED_COLLECTOR_PORT",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Enable debug messages",
			EnvVar: "FEED_COLLECTOR_DEBUG",
		},

		// Feed specific settings
		cli.StringFlag{
			Name:   "nextcloud-news-host",
			Usage:  "The hostname of the Nextclout News Server",
			EnvVar: "NEXTCLOUD_NEWS_HOST",
		},
		cli.StringFlag{
			Name:   "nextcloud-news-user",
			Usage:  "The user to use for accessing Nextcloud News",
			EnvVar: "NEXTCLOUD_NEWS_USER",
		},
		cli.StringFlag{
			Name:   "nextcloud-news-password",
			Usage:  "The password to use for accessing Nextcloud News",
			EnvVar: "NEXTCLOUD_NEWS_PASSWORD",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}

func run(c *cli.Context) {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	port := c.String("port")

	log.WithFields(log.Fields{
		"version": version,
		"port":    port,
	}).Info("Starting " + servername)

	s := Server{
		Port:                  port,
		RootEndpoint:          "/",
		NextcloudNewsHost:     c.String("nextcloud-news-host"),
		NextcloudNewsUser:     c.String("nextcloud-news-user"),
		NextcloudNewsPassword: c.String("nextcloud-news-password"),
	}
	err := s.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
