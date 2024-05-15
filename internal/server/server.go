package server

import (
	"acsm-live_timing-parser/pkg/config"
	"acsm-live_timing-parser/pkg/templating"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	router           *chi.Mux
	httpLog          *os.File
	address          string
	urlPrefix        string
	ScrapeConfig     *config.ScrapeConfig
	templatesManager *templating.TemplateManager
}

func NewServer(httpLog *os.File, address string, urlPrefix string, scrapeConfig *config.ScrapeConfig) *Server {
	return &Server{
		router:       chi.NewRouter(),
		httpLog:      httpLog,
		address:      address,
		urlPrefix:    urlPrefix,
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

	// API
	r.Post(s.urlPrefix+"/live-timing", s.Scan)

	// Web
	r.Get(s.urlPrefix+"/results", s.ResultsIndex)

	// Assets (static files)
	fs := http.FileServer(http.Dir("assets/"))
	r.Handle(s.urlPrefix+"/*", http.StripPrefix(s.urlPrefix+"/", fs))

	// Set the path for UI templates
	var err error
	s.templatesManager, err = templating.NewTemplateManager("templates/")
	if err != nil {
		slog.Error(fmt.Sprintf("Error initializing : %v", err))
		panic(err)
	}

	err = http.ListenAndServe(s.address, r)
	if err != nil {
		return err
	}

	return nil
}
