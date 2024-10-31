package server

import (
	"acsm-live_timing-parser/pkg/acsm_parser"
	"acsm-live_timing-parser/pkg/helpers"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type SetGeneralBallastRequest struct {
	Server int `json:"server"`
}

func (l *SetGeneralBallastRequest) Bind(r *http.Request) error {
	return nil
}

func (l *SetGeneralBallastRequest) ValidateRequest() error {
	return nil
}

func (s *Server) SetGeneralBallast(w http.ResponseWriter, r *http.Request) {
	defer s.Http500IfPanic(r, w)
	defer r.Body.Close()

	reqId := middleware.GetReqID(r.Context())

	// Retrieve input request data
	data := &SetGeneralBallastRequest{}

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

	// download classification list
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(s.RaceConfig.ClassificationIdsUrl)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Cannot download the classification URL (%q): %v", reqId, s.RaceConfig.ClassificationIdsUrl, err))
		render.Render(w, r, ErrInvalidRequest(err, true))
		return
	}
	defer resp.Body.Close()

	// Check if download ok
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("[%s] Cannot download the classification URL (%q): status code %d", reqId, s.RaceConfig.ClassificationIdsUrl, resp.StatusCode)
		slog.Error(msg)
		render.Render(w, r, ErrInvalidRequest(errors.New(msg), true))
		return
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Cannot read the downloaded content from %q: %v", reqId, s.RaceConfig.ClassificationIdsUrl, err))
		render.Render(w, r, ErrInvalidRequest(err, true))
		return
	}

	contentStr := string(content)

	var result []helpers.ACSMResult
	for _, record := range strings.Split(contentStr, ",") {
		result = append(result, helpers.ACSMResult{
			DriverGUID: record,
		})
	}

	jsonPath := path.Join(s.ScrapeConfig.WorkingPath, fmt.Sprintf("%d_%s", data.Server, helpers.TmpBallastFileName))
	err = acsm_parser.SaveBallast(result, jsonPath, s.RaceConfig.Ballast)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error saving ballast: %v", reqId, err))
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}
