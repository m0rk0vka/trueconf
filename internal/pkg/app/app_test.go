package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, a *App) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.r.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestRoot(t *testing.T) {
	// Create a New Server Struct
	s, _ := New()
	// Mount Handlers
	//s.MountHandlers()

	// Create a New Request
	req := httptest.NewRequest("GET", "/")
	// Execute Request
	response := executeRequest(req, s)

	// Check the response code
	checkResponseCode(t, http.StatusOK, response.Code)

	// We can use testify/require to assert values, as it is more convenient
	//require.Equal(t, "", response.Body.String())
}
