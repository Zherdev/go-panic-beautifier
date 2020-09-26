package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zherdev/go-panic-beautifier/pkg/app"
	"github.com/zherdev/go-panic-beautifier/pkg/cfg"
	"os"
)

// usage: $ go-panic-beautifier conf.yaml
func main() {
	lg := logrus.New()
	lg.SetFormatter(
		&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	)

	if len(os.Args) < 2 {
		lg.Fatal("bad cmd params, exiting")
	}

	cfgFilename := os.Args[1]
	config, err := cfg.NewConfig(cfgFilename)
	if err != nil {
		lg.WithError(err).Fatal("bad config, exiting")
	}

	a := app.NewApp(config, lg)
	a.Run()
}
