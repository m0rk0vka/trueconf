package errors

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

var USER_NOT_FOUND error = fmt.Errorf("user_not_found")

// ErrorResponse is the response that represents an error.
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Error is required by the error interface.
func (e ErrorResponse) Error() string {
	return e.Message
}

func (e ErrorResponse) StatusCode() int {
	return e.Status
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode())
	return nil
}

func BadRequest(msg string) ErrorResponse /*render.Renderer*/ {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

func NotFound(msg string) ErrorResponse /*render.Renderer*/ {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return ErrorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}

func InternalServerError(msg string) ErrorResponse /*render.Renderer*/ {
	if msg == "" {
		msg = "We encountered an error while processing your request."
	}
	return ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}
