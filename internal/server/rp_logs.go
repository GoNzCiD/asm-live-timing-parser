package server

import (
	"acsm-live_timing-parser/pkg/helpers"
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
)

func (s *Server) RpLogs(w http.ResponseWriter, r *http.Request) {
	defer s.Http500IfPanic(r, w)
	defer r.Body.Close()

	reqId := middleware.GetReqID(r.Context())

	slog.Info(fmt.Sprintf("[%s] :: Web: /rp-logs", reqId))

	pattern := r.URL.Query().Get("pattern")
	serverStr := r.URL.Query().Get("server")

	slog.Info(fmt.Sprintf("[%s] :: Web: /rp-logs {\"pattern\": %q, \"server\": %q}", reqId, pattern, serverStr))

	server, err := strconv.Atoi(serverStr)
	if err != nil {
		slog.Warn(fmt.Sprintf("[%s] Cannot parse the server parameter: %v", reqId, err))
		server = 0
	}

	files, err := helpers.FindInFolder(path.Join(s.RaceConfig.RPLogsPath, strconv.Itoa(server)), pattern)
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Error obtaining RP logs path: %v", reqId, err))
	}

	dateStr := ""
	patternArray := strings.Split(pattern, "_")
	if len(patternArray) >= 3 {
		dateStr = patternArray[2]
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"rp-logs-%s.zip\"", dateStr))
	w.WriteHeader(http.StatusOK)

	zipWriter := zip.NewWriter(w)

	for _, entry := range files {
		header := &zip.FileHeader{
			Name:     filepath.Base(entry),
			Method:   zip.Store,
			Modified: time.Now(),
		}
		entryWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			panic(err)
		}

		//fileReader := bufio.NewReader( entry)
		fileReader, err := os.Open(entry)
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] Error reading RP log file: %v", reqId, err))
			continue
		}

		_, err = io.Copy(entryWriter, fileReader)
		if err != nil {
			panic(err)
		}

		zipWriter.Flush()
		//flushingWriter, ok := z.destination.(http.Flusher)
		//if ok {
		//	flushingWriter.Flush()
		//}
	}

	err = zipWriter.Close()
	if err != nil {
		panic(err)
	}
}
