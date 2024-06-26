package server

import (
	"acsm-live_timing-parser/pkg/acsm_parser"
	"acsm-live_timing-parser/pkg/downloader"
	"acsm-live_timing-parser/pkg/helpers"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type LiveTimingRequest struct {
	Server         int    `json:"server"`
	ForceDownload  bool   `json:"force_download"`
	Name           string `json:"name"`
	PreviewPattern string `json:"preview_pattern"`
}

type LiveTimingResult struct {
	Success        bool             `json:"-"`
	HTTPStatusCode int              `json:"-"`
	Error          string           `json:"error_message,omitempty"`
	Data           []helpers.Hotlap `json:"data,omitempty"`
}

func (l *LiveTimingRequest) Bind(r *http.Request) error {
	return nil
}

func (l *LiveTimingRequest) ValidateRequest() error {
	if l.Name == "" {
		return fmt.Errorf("'Name' parameter is required")
	}
	return nil
}

func (s *Server) Scan(w http.ResponseWriter, r *http.Request) {
	defer s.Http500IfPanic(r, w)
	defer r.Body.Close()

	reqId := middleware.GetReqID(r.Context())

	//if !s.IsAccessGranted(logger, w, r) {
	//	return
	//}

	// Retrieve input request data
	data := &LiveTimingRequest{}

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

	// Log input data to process
	slog.Info(fmt.Sprintf("[%s] :: Endpoint: /live-timing: %+v", reqId, data))

	end := make(chan LiveTimingResult, 1)
	go func() {
		var jsonStr string
		var downloaded bool
		var liveTiming *helpers.LeaderBoard
		jsonPath := path.Join(s.ScrapeConfig.WorkingPath, fmt.Sprintf("%d_%s", data.Server, helpers.TmpLeaderBoardFileName))
		info, err := os.Stat(jsonPath)
		if err != nil {
			slog.Debug(fmt.Sprintf("[%s] Cannot retrieve the temporal leaderboard json %q: %v", reqId, jsonPath, err))
		}
		if data.ForceDownload || err != nil || time.Since(info.ModTime()) > s.ScrapeConfig.LeaderBoardJsonTtl {
			slog.Debug(fmt.Sprintf("[%s] Retrieving leaderboard for server %d", reqId, data.Server))
			downloadHandler := downloader.NewDownloader(
				s.ScrapeConfig.ChromeDriverPath,
				s.ScrapeConfig.SeleniumUrl,
				s.ScrapeConfig.User,
				s.ScrapeConfig.Password,
				s.ScrapeConfig.ACSMDomain,
				data.Server,
				false)
			var dErr error
			jsonStr, dErr = downloadHandler.Download(downloader.LiveTimingApiEndpoint, nil)
			if dErr != nil {
				slog.Error(fmt.Sprintf("[%s] Cannot download leaderboard json: %v", reqId, dErr))
				// In case of tmp file doesn't exist return with error
				if err != nil {
					end <- LiveTimingResult{
						Success:        false,
						HTTPStatusCode: 500,
						Error:          dErr.Error(),
					}
					return
				}
			}

			if jsonStr == "Too Many Requests" {
				slog.Error(fmt.Sprintf("[%s] Too many request to API", reqId))
				end <- LiveTimingResult{
					Success:        false,
					HTTPStatusCode: 429,
					Error:          "Too Many Request",
				}
				return
			}

			slog.Debug(fmt.Sprintf("[%s] Parsing JSON retrieved from API", reqId))
			liveTiming, err = acsm_parser.ReadLeaderBoardJson(jsonStr)
			if err != nil {
				slog.Error(fmt.Sprintf("[%s] Cannot parse leaderboard json value %q: %v", reqId, jsonPath, err))
				end <- LiveTimingResult{
					Success:        false,
					HTTPStatusCode: 500,
					Error:          err.Error(),
				}
				return
			}

			if liveTiming.Name == data.Name {
				downloaded = true
			} else {
				slog.Info(fmt.Sprintf("[%s] Download content does not match with input name parameter %q", reqId, data.Name))
			}
		}

		if !downloaded { //jsonStr == "" {
			slog.Debug(fmt.Sprintf("[%s] Reading JSON retrieved from temporal file", reqId))
			jsonStr, err = helpers.ReadFromFile(jsonPath)
			if err != nil {
				slog.Error(fmt.Sprintf("[%s] Cannot read temporal leaderboard json %q: %v", reqId, jsonPath, err))
				end <- LiveTimingResult{
					Success:        false,
					HTTPStatusCode: 500,
					Error:          err.Error(),
				}
				return
			}
			if jsonStr != "" {
				slog.Debug(fmt.Sprintf("[%s] Parsing JSON retrieved from temporal file", reqId))
				liveTiming, err = acsm_parser.ReadLeaderBoardJson(jsonStr)
				if err != nil {
					slog.Error(fmt.Sprintf("[%s] Cannot parse leaderboard json value %q: %v", reqId, jsonPath, err))
					end <- LiveTimingResult{
						Success:        false,
						HTTPStatusCode: 500,
						Error:          err.Error(),
					}
					return
				}
			}
		}

		if jsonStr == "" || (liveTiming != nil && liveTiming.Name != data.Name) {
			end <- LiveTimingResult{
				Success:        false,
				HTTPStatusCode: 404,
			}
			return
		}

		slog.Debug(fmt.Sprintf("[%s] Reading ballast", reqId))
		ballastPath := path.Join(s.ScrapeConfig.WorkingPath, fmt.Sprintf("%d_%s", data.Server, helpers.TmpBallastFileName))
		ballast, err := acsm_parser.ReadAndParseBallast(ballastPath)
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] Cannot retrieve ballast from %q: %v", reqId, ballastPath, err))
		}

		slog.Debug(fmt.Sprintf("[%s] Extracting hotlaps", reqId))
		result, bestSectors := acsm_parser.ExtractHotlaps(append(liveTiming.ConnectedDrivers, liveTiming.DisconnectedDrivers...))
		result = acsm_parser.SortHotlapsAndCalculateData(result, &data.PreviewPattern, bestSectors, ballast)

		if downloaded && jsonStr != "" {
			slog.Debug(fmt.Sprintf("[%s] Saving leaderboard json to temp folder %q", reqId, jsonPath))
			err = helpers.SaveToFile(jsonStr, jsonPath)
			if err != nil {
				slog.Error(fmt.Sprintf("[%s] Cannot save temporal leaderboard json %q: %v", reqId, jsonPath, err))
			}
		}

		end <- LiveTimingResult{
			Success:        true,
			HTTPStatusCode: 200,
			Data:           result,
		}
	}()

	select {
	case <-r.Context().Done():
		slog.Error(fmt.Sprintf("[%s] Timeout", reqId))
		res := NewTimeoutResponse()
		res.RequestID = reqId
		w.WriteHeader(res.HTTPStatusCode)
		json.NewEncoder(w).Encode(res)
	case r := <-end:
		slog.Info(fmt.Sprintf("[%s] Completed", reqId))
		w.WriteHeader(r.HTTPStatusCode)
		if r.Error != "" {
			slog.Info(fmt.Sprintf("[%s] Some error happen: %s", reqId, r.Error))
			json.NewEncoder(w).Encode(r)
		} else {
			slog.Debug(fmt.Sprintf("[%s] Retrieved data: %+v", reqId, r.Data))
			json.NewEncoder(w).Encode(r.Data)
		}
	}
}
