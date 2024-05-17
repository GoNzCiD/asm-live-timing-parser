package server

import (
	"acsm-live_timing-parser/pkg/acsm_parser"
	"acsm-live_timing-parser/pkg/downloader"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
)

type ResultsListData struct {
	Track        string
	Layout       string
	Type         string
	Date         time.Time
	DownloadLink string
	RPPattern    string
}

func extractRPPattern(jsonPath string) string {
	//jsonPath = "results/download/2024_5_9_19_48_RACE.json"
	a := strings.Split(jsonPath, "/")
	a = strings.Split(a[len(a)-1], "_")

	year := a[0]
	month := a[1]
	if len(month) == 1 {
		month = "0" + month
	}
	day := a[2]
	if len(day) == 1 {
		day = "0" + day
	}

	return fmt.Sprintf("race_penalty_%s%s%s_*.log", year, month, day)
}

func (s *Server) ResultsIndex(w http.ResponseWriter, r *http.Request) {
	defer s.Http500IfPanic(r, w)
	defer r.Body.Close()

	reqId := middleware.GetReqID(r.Context())

	slog.Info(fmt.Sprintf("[%s] :: Web: /results", reqId))

	data := struct {
		Title    string
		ServerNo int
		Results  []ResultsListData
	}{
		Title:    "Results list",
		ServerNo: 2,
	}

	downloadHandler := downloader.NewDownloader(
		s.ScrapeConfig.ChromeDriverPath,
		s.ScrapeConfig.SeleniumUrl,
		s.ScrapeConfig.User,
		s.ScrapeConfig.Password,
		s.ScrapeConfig.ACSMDomain,
		data.ServerNo,
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
			RPPattern:    extractRPPattern(r.ResultsJSONURL),
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

func (s *Server) Results(w http.ResponseWriter, r *http.Request) {
	defer s.Http500IfPanic(r, w)
	defer r.Body.Close()

	reqId := middleware.GetReqID(r.Context())

	resultPath := r.URL.Path[7:]
	serverStr := r.URL.Query().Get("server")

	slog.Info(fmt.Sprintf("[%s] :: Web: /result {\"path\": %q, \"server\": %q}", reqId, resultPath, serverStr))

	server, err := strconv.Atoi(serverStr)
	if err != nil {
		slog.Warn(fmt.Sprintf("[%s] Cannot parse the server parameter: %v", reqId, err))
		server = 0
	}

	downloadHandler := downloader.NewDownloader(
		s.ScrapeConfig.ChromeDriverPath,
		s.ScrapeConfig.SeleniumUrl,
		s.ScrapeConfig.User,
		s.ScrapeConfig.Password,
		s.ScrapeConfig.ACSMDomain,
		server,
		false)
	jsonStr, err := downloadHandler.Download(resultPath, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error retrieving results: %v", reqId, err))
	}

	results, err := acsm_parser.GetResults(jsonStr)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error parsing results: %v", reqId, err))
	}

	name := filepath.Base(resultPath)
	name = strings.Split(name, ".")[0]

	file, err := acsm_parser.GenerateResultsExcel(results, name, s.RaceConfig.Points)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error generating excel results: %v", reqId, err))
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", name+".xlsx"))
	err = file.Write(w)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error exporting excel: %v", reqId, err))
	}
}
