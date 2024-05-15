package server

import (
	"acsm-live_timing-parser/pkg/acsm_parser"
	"acsm-live_timing-parser/pkg/downloader"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

type ResultsListData struct {
	Track        string
	Layout       string
	Type         string
	Date         time.Time
	DownloadLink string
}

func (s *Server) ResultsIndex(w http.ResponseWriter, r *http.Request) {
	defer s.Http500IfPanic(r, w)
	defer r.Body.Close()

	reqId := middleware.GetReqID(r.Context())

	slog.Info(fmt.Sprintf("[%s] :: Web: /results", reqId))

	data := struct {
		Title   string
		Results []ResultsListData
	}{
		Title: "Results list",
	}

	downloadHandler := downloader.NewDownloader(
		s.ScrapeConfig.ChromeDriverPath,
		s.ScrapeConfig.SeleniumUrl,
		s.ScrapeConfig.User,
		s.ScrapeConfig.Password,
		s.ScrapeConfig.ACSMDomain,
		2,
		false)
	jsonStr, err := downloadHandler.Download(downloader.ResultsListApiEndpoint, map[string]string{"q": "Type:\"RACE\""})
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error retrieving results list: %v", reqId, err))
	}

	results, err := acsm_parser.ReadResultsJson(jsonStr)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error retrieving results list: %v", reqId, err))
	}
	for _, r := range results.Results {
		data.Results = append(data.Results, ResultsListData{
			Track:        r.Track,
			Layout:       r.TrackLayout,
			Type:         r.SessionType,
			Date:         r.Date,
			DownloadLink: r.ResultsJSONURL,
		})
	}

	t, err := s.templatesManager.Template("results/index")
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error retrieving template: %v", reqId, err))
	}
	err = t.Execute(w, data)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error rendering template: %v", reqId, err))
	}
}
