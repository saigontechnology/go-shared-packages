package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/datngo2sgtech/go-packages/api"
)

// CT is a Component Test utility/**
// This package is used to implement component test for endpoints.
type CT interface {
	IT
	API() api.API
	// SetApi is used to set a concrete api of service
	SetAPI(api api.API)
	Request() *http.Request
	SetRequest(req *http.Request)
	Recorder() *httptest.ResponseRecorder
	// AssertRestResponse is used to verify response of a REST endpoint
	AssertRestResponse(statusCode int, responseBody string)
	// RequireRestResponse is used to verify response of a REST endpoint
	RequireRestResponse(statusCode int, responseBody string)
}

// SqliteCt is a Component Test utility with Sqlite database.
type SqliteCt struct {
	*it
	api      api.API
	request  *http.Request
	recorder *httptest.ResponseRecorder
}

func NewSqliteCt(t *testing.T, options ...ItOption) *SqliteCt {
	t.Helper()

	it := newIT(t, options...)
	return &SqliteCt{
		it:       it,
		recorder: httptest.NewRecorder(),
	}
}

func (t *SqliteCt) API() api.API {
	return t.api
}

func (t *SqliteCt) SetAPI(api api.API) {
	t.api = api
}

func (t *SqliteCt) Request() *http.Request {
	return t.request
}

func (t *SqliteCt) SetRequest(req *http.Request) {
	t.request = req
}

func (t *SqliteCt) Recorder() *httptest.ResponseRecorder {
	return t.recorder
}

func (t *SqliteCt) AssertRestResponse(
	statusCode int,
	responseBody string,
) {
	t.assert.Equal(statusCode, t.recorder.Code)
	t.assert.Equal(responseBody, t.recorder.Body.String())
}

func (t *SqliteCt) RequireRestResponse(
	statusCode int,
	responseBody string,
) {
	responseBody = strings.ReplaceAll(responseBody, "\n", "")
	responseBody = strings.ReplaceAll(responseBody, "\t", "")
	responseBody = strings.TrimSpace(responseBody)

	t.require.Equal(statusCode, t.recorder.Code)
	t.require.Equal(responseBody, t.recorder.Body.String())
}
