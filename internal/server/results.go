package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
)

func (s *Server) ResultsIndex(w http.ResponseWriter, r *http.Request) {
	defer s.Http500IfPanic(r, w)
	defer r.Body.Close()

	reqId := middleware.GetReqID(r.Context())

	slog.Info(fmt.Sprintf("[%s] :: Web: /results", reqId))

	data := struct {
		Title string
	}{
		Title: "Results list",
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
