package server

import (
	"acsm-live_timing-parser/pkg/config"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	router       *chi.Mux
	httpLog      *os.File
	address      string
	ScrapeConfig *config.ScrapeConfig
}

func NewServer(httpLog *os.File, address string, scrapeConfig *config.ScrapeConfig) *Server {
	return &Server{
		router:       chi.NewRouter(),
		httpLog:      httpLog,
		address:      address,
		ScrapeConfig: scrapeConfig,
	}
}

func (s *Server) InitializeAndStart() error {
	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{Logger: log.New(s.httpLog, "", log.LstdFlags), NoColor: true})

	r := s.router
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Post("/live-timing", s.Scan)

	err := http.ListenAndServe(s.address, r)
	if err != nil {
		return err
	}

	return nil
}
