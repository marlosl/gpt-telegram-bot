package tests

import (
	"net/http"
	"net/http/httptest"
)

func CreateTestServer(fn func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(fn))
}
