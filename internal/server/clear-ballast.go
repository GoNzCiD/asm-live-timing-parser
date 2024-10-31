package server

import (
	"acsm-live_timing-parser/pkg/helpers"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type ClearBallastRequest struct {
	Server int `json:"server"`
}

func (l *ClearBallastRequest) Bind(r *http.Request) error {
	return nil
}

func (l *ClearBallastRequest) ValidateRequest() error {
	return nil
}

func (s *Server) ClearBallast(w http.ResponseWriter, r *http.Request) {
	defer s.Http500IfPanic(r, w)
	defer r.Body.Close()

	reqId := middleware.GetReqID(r.Context())

	// Retrieve input request data
	data := &ClearBallastRequest{}

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

	jsonPath := path.Join(s.ScrapeConfig.WorkingPath, fmt.Sprintf("%d_%s", data.Server, helpers.TmpBallastFileName))
	if helpers.CheckFileExists(jsonPath) {
		err := os.Remove(jsonPath)
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] Error clearing ballast: %v", reqId, err))
			panic(err)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}
