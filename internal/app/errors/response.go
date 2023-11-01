package errors

import (
	"net/http"

	"github.com/go-chi/render"
)

const USER_NOT_FOUND string = "User not found"

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

func BadRequest(msg string) render.Renderer {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return &ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

func NotFound(msg string) render.Renderer {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return &ErrorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}
