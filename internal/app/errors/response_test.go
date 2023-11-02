package errors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	assert "github.com/go-playground/assert/v2"
)

func TestErrorResponse_Error(t *testing.T) {
	e := ErrorResponse{
		Message: "abc",
	}
	assert.Equal(t, "abc", e.Error())
}

func TestErrorResponse_StatusCode(t *testing.T) {
	e := ErrorResponse{
		Status: 400,
	}
	assert.Equal(t, 400, e.StatusCode())
}

func TestErrorResponse_Render(t *testing.T) {
	e := ErrorResponse{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	assert.Equal(t, nil, e.Render(w, r))
}

func TestInternalServerError(t *testing.T) {
	msg := "test"
	res := InternalServerError(msg)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
	assert.Equal(t, msg, res.Error())
}

func TestNotFound(t *testing.T) {
	msg := "test"
	res := NotFound(msg)
	assert.Equal(t, http.StatusNotFound, res.StatusCode())
	assert.Equal(t, msg, res.Error())
}

func TestBadRequest(t *testing.T) {
	msg := "test"
	res := BadRequest(msg)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	assert.Equal(t, msg, res.Error())
}
