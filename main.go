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
	app.HideHelp = false

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "Show this help message",
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "Show the version",
	}
	cli.AppHelpTemplate = `NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
`

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
			Usage:  "The hostname of the Nextcloud News Server",
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
		cli.StringFlag{
			Name:   "pocket-access-token",
			Usage:  "The token for a Pocket user",
			EnvVar: "POCKET_ACCESS_TOKEN",
		},
		cli.StringFlag{
			Name:   "pocket-consumer-key",
			Usage:  "The key for accessing the Pocket API",
			EnvVar: "POCKET_CONSUMER_KEY",
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
		PocketAccessToken:     c.String("pocket-access-token"),
		PocketConsumerKey:     c.String("pocket-consumer-key"),
	}
	err := s.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
