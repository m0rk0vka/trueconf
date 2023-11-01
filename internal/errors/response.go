package errors

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var UserNotFound = errors.New("user_not_found")

// ErrResponse is the response that represents an error.
type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func BadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Your request is in a bad format.",
		ErrorText:      err.Error(),
	}
}
