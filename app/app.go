package app

import (
	"context"
	"log/slog"
	"net/http"
)

type App struct {
	Name   string
	Addr   string
	Logger *slog.Logger
	Server *http.Server
}

func New(name string, addr string, handler http.Handler, logger *slog.Logger) *App {
	return &App{
		Name:   name,
		Addr:   addr,
		Logger: logger,
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (a *App) Run() error {
	if a.Logger != nil {
		a.Logger.Info("application starting", "name", a.Name, "addr", a.Addr)
	}
	return a.Server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.Logger != nil {
		a.Logger.Info("application shutting down", "name", a.Name)
	}
	return a.Server.Shutdown(ctx)
}
