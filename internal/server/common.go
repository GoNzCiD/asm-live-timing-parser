package server

import (
	"acsm-live_timing-parser/pkg/helpers"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

// Http500IfPanic checks whether there panic pending to be caught.
func (s *Server) Http500IfPanic(r interface{}, w http.ResponseWriter) {
	err := recover()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Unexpected error"
		_, _ = w.Write([]byte(msg))

		slog.Error(fmt.Sprintf("unexpected panic: %v", err))
		slog.Error(fmt.Sprintf("request: %v", r))
		slog.Error(fmt.Sprintf("stack trace:"))
		for _, s := range helpers.StackTraceAsList() {
			slog.Error(fmt.Sprintf("%s", s))
		}
	}
}

// ErrorResponse renderer type for handling all sorts of errors.
type ErrorResponse struct {
	Err            error  `json:"-"` // low-level runtime error
	HTTPStatusCode int    `json:"-"` // http response status code
	Error          string `json:"error_message"`
	RequestID      string `json:"request_id,omitempty"`
}

// Error response render
func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrResponse - returns an HTTP response format when an internal error occurs
func ErrResponse(err error, custom bool) render.Renderer {
	errorMessage := "Error processing the request."
	if custom {
		errorMessage = err.Error()
	}
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: 500,
		Error:          errorMessage,
	}
}

// ErrInvalidRequest - returns an HTTP response format when an error occurs
// while parsing & validating input request
// custom - when true shows the concrete error
func ErrInvalidRequest(err error, custom bool) render.Renderer {
	errorMessage := "Invalid input structure."
	if custom {
		errorMessage = err.Error()
	}
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: 400,
		Error:          errorMessage,
	}
}

func NewTimeoutResponse() *ErrorResponse {
	result := &ErrorResponse{
		Err:            nil,
		HTTPStatusCode: 408,
		Error:          "Timeout",
	}

	return result
}
