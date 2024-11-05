package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/5aradise/link-forge/config"
	"github.com/5aradise/link-forge/internal/database"
	"github.com/5aradise/link-forge/internal/handlers"
	"github.com/5aradise/link-forge/pkg/httpserver"
	"github.com/5aradise/link-forge/pkg/logger"
	"github.com/5aradise/link-forge/pkg/middleware"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	// _ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load config
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Create logger
	l := logger.New(os.Stdout, config.Cfg.Env)

	// Connect to storage
	conn, err := sql.Open("libsql", config.Cfg.DB.URL) // sqlite3
	if err != nil {
		l.Error("sql open", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close()

	db := database.Create(conn)

	// Set handlers
	router := http.NewServeMux()
	router.HandleFunc("/healthz", handlers.Readiness(l))

	// api
	api := http.NewServeMux()

	// v1
	v1 := http.NewServeMux()

	URLService := handlers.NewURLService(l, db)
	v1.HandleFunc(http.MethodPost+" /urls", URLService.CreateURL)
	v1.HandleFunc(http.MethodGet+" /urls", URLService.ListURLs)
	v1.HandleFunc(http.MethodGet+" /urls/{alias}", URLService.RedirectURL)
	v1.HandleFunc(http.MethodDelete+" /urls/{alias}", URLService.DeleteURL)

	api.Handle("/v1/", http.StripPrefix("/v1", v1))

	router.Handle("/api/", http.StripPrefix("/api", api))

	// Run server
	server := httpserver.New(
		middleware.Use(router,
			middleware.Recoverer(l),
			middleware.Cors(l),
			middleware.RequestID(l),
			middleware.Logger(l),
		),
		httpserver.Port(config.Cfg.Server.Port),
		httpserver.ReadTimeout(config.Cfg.Server.Timeout),
		httpserver.IdleTimeout(config.Cfg.Server.IdleTimeout),
		httpserver.ErrorLog(slog.NewLogLogger(l.With(slog.String("source", "httpserver")).Handler(), slog.LevelError)),
	)

	l.Info("starting server", slog.String("address", server.Addr()))
	go server.Run()

	// Waiting signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Error("signal interrupt", slog.String("error", s.String()))
	case err := <-server.Notify():
		l.Error("server notify", slog.String("error", err.Error()))
	}

	// Shutdown server
	err = server.Shutdown()
	if err != nil {
		l.Info("server shutdown", slog.String("error", err.Error()))
	}
}
