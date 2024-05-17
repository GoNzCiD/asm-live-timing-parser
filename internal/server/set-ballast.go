package server

import (
	"acsm-live_timing-parser/pkg/acsm_parser"
	"acsm-live_timing-parser/pkg/downloader"
	"acsm-live_timing-parser/pkg/helpers"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type SetBallastRequest struct {
	Server      int    `json:"server"`
	ResultsPath string `json:"results"`
}

func (l *SetBallastRequest) Bind(r *http.Request) error {
	return nil
}

func (l *SetBallastRequest) ValidateRequest() error {
	if l.ResultsPath == "" {
		return fmt.Errorf("'results' parameter is required")
	}
	return nil
}

func (s *Server) SetBallast(w http.ResponseWriter, r *http.Request) {
	defer s.Http500IfPanic(r, w)
	defer r.Body.Close()

	reqId := middleware.GetReqID(r.Context())

	// Retrieve input request data
	data := &SetBallastRequest{}

	// check the overall shape of the request body
	if err := render.Bind(r, data); err != nil {
		slog.Error(fmt.Sprintf("[%s] Invalid data provided: %v", reqId, err))
		render.Render(w, r, ErrInvalidRequest(err, false))
		return
	}

	// validate the content of the parsed request body
	if err := data.ValidateRequest(); err != nil {
		slog.Error(fmt.Sprintf("[%s] Invalid data provided: %v", reqId, err))
		render.Render(w, r, ErrInvalidRequest(err, true))
		return
	}

	downloadHandler := downloader.NewDownloader(
		s.ScrapeConfig.ChromeDriverPath,
		s.ScrapeConfig.SeleniumUrl,
		s.ScrapeConfig.User,
		s.ScrapeConfig.Password,
		s.ScrapeConfig.ACSMDomain,
		data.Server,
		false)
	jsonStr, err := downloadHandler.Download(data.ResultsPath, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error retrieving results: %v", reqId, err))
		panic(err)
	}

	results, err := acsm_parser.GetResults(jsonStr)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error parsing results: %v", reqId, err))
		panic(err)
	}

	jsonPath := path.Join(s.ScrapeConfig.WorkingPath, fmt.Sprintf("%d_%s", data.Server, helpers.TmpBallastFileName))
	err = acsm_parser.SaveBallast(results.Result, jsonPath, s.RaceConfig.Ballast)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error saving ballast: %v", reqId, err))
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}
