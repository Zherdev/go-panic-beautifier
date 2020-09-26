package app

import (
	"github.com/sirupsen/logrus"
	"github.com/zherdev/go-panic-beautifier/pkg/cfg"
	"github.com/zherdev/go-panic-beautifier/pkg/handlers"
	"net/http"
)

type App struct {
	config *cfg.Config
	lg     *logrus.Logger
}

func NewApp(config *cfg.Config, lg *logrus.Logger) *App {
	a := &App{config: config}

	a.setLoggerFrom(lg)

	return &App{
		config: config,
		lg:     lg,
	}
}

func (a *App) setLoggerFrom(lg *logrus.Logger) {
	a.lg = lg

	if a.config == nil {
		return
	}

	a.lg.SetLevel(logrus.Level(a.config.LogLevel))
}

func (a *App) Run() {
	a.lg.WithField("configuration", a.config).Info("Starting application")

	h, err := handlers.NewMainPageHandler(a.lg, a.config)
	if err != nil {
		a.lg.WithError(err).Fatal("fatal error at starting app")
	}

	mux := http.NewServeMux()

	mux.Handle(a.config.MainPageUrl, h)
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	)

	err = a.listenAndServe(mux)

	a.lg.Info("Stopping application")
	if err != nil {
		a.lg.WithError(err).Fatal("fatal error")
	}
}

func (a *App) listenAndServe(mux *http.ServeMux) error {
	server := &http.Server{}

	server.Addr = a.config.Host + ":" + a.config.Port
	server.ReadHeaderTimeout = a.config.ReadTimeout
	server.WriteTimeout = a.config.WriteTimout
	server.Handler = mux

	return server.ListenAndServe()
}
