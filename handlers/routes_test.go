package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mtlynch/logpaste/store/test_sqlite"
)

func TestRoutes(t *testing.T) {
	t.Chdir("..")
	ds := test_sqlite.New()
	if err := ds.InsertEntry("12345678", "dummy entry"); err != nil {
		t.Fatal(err)
	}
	s := defaultServer{
		router:       http.NewServeMux(),
		store:        ds,
		maxCharLimit: MaxPasteCharacters,
	}
	s.routes()

	for _, tt := range []struct {
		description   string
		method        string
		path          string
		status        int
		body          string
		expectedAllow string
		expectedCORS  string
	}{
		{
			description: "GET root",
			method:      http.MethodGet,
			path:        "/",
			status:      http.StatusOK,
			body:        "<!doctype html>",
		},
		{
			description: "HEAD root",
			method:      http.MethodHead,
			path:        "/",
			status:      http.StatusOK,
		},
		{
			description: "GET paste",
			method:      http.MethodGet,
			path:        "/12345678",
			status:      http.StatusOK,
			body:        "dummy entry",
		},
		{
			description: "HEAD paste",
			method:      http.MethodHead,
			path:        "/12345678",
			status:      http.StatusOK,
		},
		{
			description: "GET paste with extra segment",
			method:      http.MethodGet,
			path:        "/12345678/extra",
			status:      http.StatusNotFound,
		},
		{
			description: "GET unknown path with extra segment",
			method:      http.MethodGet,
			path:        "/unknown/extra",
			status:      http.StatusNotFound,
		},
		{
			description:   "PUT single segment",
			method:        http.MethodPut,
			path:          "/extra",
			status:        http.StatusMethodNotAllowed,
			expectedAllow: "GET, HEAD",
		},
		{
			description:   "POST single segment",
			method:        http.MethodPost,
			path:          "/extra",
			status:        http.StatusMethodNotAllowed,
			expectedAllow: "GET, HEAD",
		},
		{
			description:   "OPTIONS single segment",
			method:        http.MethodOptions,
			path:          "/extra",
			status:        http.StatusMethodNotAllowed,
			expectedAllow: "GET, HEAD",
		},
		{
			description: "PUT extra segment",
			method:      http.MethodPut,
			path:        "/12345678/extra",
			status:      http.StatusNotFound,
		},
		{
			description: "POST extra segment",
			method:      http.MethodPost,
			path:        "/12345678/extra",
			status:      http.StatusNotFound,
		},
		{
			description:   "DELETE root",
			method:        http.MethodDelete,
			path:          "/",
			status:        http.StatusMethodNotAllowed,
			expectedAllow: "GET, HEAD, OPTIONS, POST, PUT",
		},
		{
			description:  "OPTIONS root",
			method:       http.MethodOptions,
			path:         "/",
			status:       http.StatusOK,
			expectedCORS: "*",
		},
		{
			description: "GET favicon",
			method:      http.MethodGet,
			path:        "/favicon.ico",
			status:      http.StatusNotFound,
		},
		{
			description: "GET existing static file",
			method:      http.MethodGet,
			path:        "/static/css/style.css",
			status:      http.StatusOK,
		},
		{
			description: "GET missing static file",
			method:      http.MethodGet,
			path:        "/static/not-found.css",
			status:      http.StatusNotFound,
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			s.router.ServeHTTP(rr, req)
			if rr.Code != tt.status {
				t.Fatalf("status = %d, want %d", rr.Code, tt.status)
			}
			if tt.body != "" && !strings.Contains(rr.Body.String(), tt.body) {
				t.Errorf("body does not contain %q", tt.body)
			}
			if got, want := rr.Header().Get("Allow"), tt.expectedAllow; got != want {
				t.Errorf("Allow header = %q, want %q", got, want)
			}
			if got, want := rr.Header().Get("Access-Control-Allow-Origin"), tt.expectedCORS; got != want {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, want)
			}
		})
	}
}
