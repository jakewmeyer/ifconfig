package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var s *server

func TestMain(m *testing.M) {
	s, _ = new(":7000")
	code := m.Run()
	os.Exit(code)
}

// executeRequest takes a request and records the response
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)
	return rr
}

func TestReturnPlaintext(t *testing.T) {
	t.Parallel()

	s.Router.Get("/", parseIp(s))
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124")
	res := executeRequest(req)

	assert.Equal(t, http.StatusOK, res.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", res.Body.String(), "should return an ip address")
}

func TestParseJson(t *testing.T) {
	t.Parallel()

	s.Router.Get("/", parseIp(s))
	req, _ := http.NewRequest("GET", "/?json", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124")
	res := executeRequest(req)

	assert.Equal(t, http.StatusOK, res.Code, "should return a 200 status")
	assert.Equal(t, "{\"ip\":\"192.168.1.124\"}\n", res.Body.String(), "should return an ip address in json format")
}

func TestParseNoHeader(t *testing.T) {
	t.Parallel()

	s.Router.Get("/", parseIp(s))
	req, _ := http.NewRequest("GET", "/", nil)
	res := executeRequest(req)

	assert.Equal(t, http.StatusInternalServerError, res.Code, "should return a 500")
}

func TestParseMultipleIp(t *testing.T) {
	t.Parallel()

	s.Router.Get("/", parseIp(s))
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124,10.0.0.1")
	res := executeRequest(req)

	assert.Equal(t, http.StatusOK, res.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", res.Body.String(), "should return an ip address")
}

func TestParseNoIp(t *testing.T) {
	t.Parallel()

	s.Router.Get("/", parseIp(s))
	req, _ := http.NewRequest("GET", "/", nil)
	res := executeRequest(req)

	assert.Equal(t, http.StatusInternalServerError, res.Code, "should return a 500")
}
