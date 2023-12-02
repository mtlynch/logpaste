package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/logpaste/store"
)

type mockStore struct {
	entries map[string]string
}

func (ds mockStore) GetEntry(id string) (string, error) {
	if contents, ok := ds.entries[id]; ok {
		return contents, nil
	}
	return "", store.EntryNotFoundError{ID: id}
}

func (ds *mockStore) InsertEntry(id string, contents string) error {
	ds.entries[id] = contents
	return nil
}

func (ds *mockStore) Reset() {
	ds.entries = make(map[string]string)
}

func TestPasteGet(t *testing.T) {
	ds := mockStore{
		entries: map[string]string{
			"12345678": "dummy entry",
		},
	}
	router := mux.NewRouter()
	s := defaultServer{
		store:  &ds,
		router: router,
	}
	s.routes()

	for _, tt := range []struct {
		description     string
		id              string
		statusExpected  int
		contentExpected string
	}{
		{
			description:     "valid entry",
			id:              "12345678",
			statusExpected:  http.StatusOK,
			contentExpected: "dummy entry",
		},
		{
			description:     "non-existent entry",
			id:              "missing1",
			statusExpected:  http.StatusNotFound,
			contentExpected: "",
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			if got, want := w.Code, tt.statusExpected; got != want {
				t.Fatalf("status=%d, want=%d", got, want)
			}
			if w.Code != http.StatusOK {
				return
			}
			bodyBytes, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatalf("failed to read HTTP response body: %v", err)
			}
			if got, want := string(bodyBytes), tt.contentExpected; got != want {
				t.Errorf("body=%s, want=%s", got, want)
			}
		})
	}
}

func TestPastePut(t *testing.T) {
	for _, tt := range []struct {
		description    string
		body           string
		statusExpected int
	}{
		{
			description:    "valid content",
			body:           "hello, world!",
			statusExpected: http.StatusOK,
		},
		{
			description:    "just at size limit",
			body:           strings.Repeat("A", MaxPasteCharacters),
			statusExpected: http.StatusOK,
		},
		{
			description:    "too long content",
			body:           strings.Repeat("A", MaxPasteCharacters+1),
			statusExpected: http.StatusBadRequest,
		},
		{
			description:    "empty content",
			body:           "",
			statusExpected: http.StatusBadRequest,
		},
		{
			description:    "just whitespace",
			body:           "  ",
			statusExpected: http.StatusBadRequest,
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			ds := mockStore{
				entries: make(map[string]string),
			}
			router := mux.NewRouter()
			s := defaultServer{
				store:  &ds,
				router: router,
			}
			s.routes()

			ds.Reset()

			req, err := http.NewRequest("PUT", "/", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			if got, want := w.Code, tt.statusExpected; got != want {
				t.Fatalf("status=%d, want=%d", got, want)
			}
			if w.Code != http.StatusOK {
				return
			}

			var resp PastePutResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to read HTTP response body: %v", err)
			}
			storedContents, err := ds.GetEntry(resp.Id)
			if err != nil {
				t.Fatalf("no entry found with id: %s", resp.Id)
			}
			if got, want := storedContents, tt.body; got != want {
				t.Fatalf("stored=%v, want=%v", got, want)
			}
		})
	}
}

func TestPastePost(t *testing.T) {
	ds := mockStore{
		entries: make(map[string]string),
	}
	router := mux.NewRouter()
	s := defaultServer{
		store:  &ds,
		router: router,
	}
	s.routes()
	for _, tt := range []struct {
		description      string
		contentType      string
		body             string
		statusExpected   int
		contentsExpected string
	}{
		{
			description:    "reject non-multipart data",
			contentType:    "text/plain",
			body:           "hello, world!",
			statusExpected: http.StatusBadRequest,
		},
		{
			description: "rejects empty input",
			contentType: "multipart/form-data; boundary=------------------------aea33768a2527972",
			body: `
--------------------------aea33768a2527972
Content-Disposition: form-data; name="dummyname1"



--------------------------aea33768a2527972--`,
			statusExpected: http.StatusBadRequest,
		},
		{
			description: "accepts string data",
			contentType: "multipart/form-data; boundary=------------------------aea33768a2527972",
			body: `
--------------------------aea33768a2527972
Content-Disposition: form-data; name="dummyname2"

some data I want to upload
--------------------------aea33768a2527972--`,
			statusExpected:   http.StatusOK,
			contentsExpected: "some data I want to upload",
		},
		{
			description: "accepts string data at size limit",
			contentType: "multipart/form-data; boundary=------------------------aea33768a2527972",
			body: fmt.Sprintf(`
--------------------------aea33768a2527972
Content-Disposition: form-data; name="dummyname2"

%s
--------------------------aea33768a2527972--`, strings.Repeat("A", MaxPasteCharacters)),
			statusExpected:   http.StatusOK,
			contentsExpected: strings.Repeat("A", MaxPasteCharacters),
		},
		{
			description: "rejects string data above size limit",
			contentType: "multipart/form-data; boundary=------------------------aea33768a2527972",
			body: fmt.Sprintf(`
--------------------------aea33768a2527972
Content-Disposition: form-data; name="dummyname2"

%s
--------------------------aea33768a2527972--`, strings.Repeat("A", MaxPasteCharacters+1)),
			statusExpected: http.StatusBadRequest,
		},
		{
			description: "accepts file data",
			contentType: "multipart/form-data; boundary=------------------------ff01448fc0d75457",
			body: `
--------------------------ff01448fc0d75457
Content-Disposition: form-data; name="dummyname3"; filename="text.txt"
Content-Type: text/plain

some data in a file
--------------------------ff01448fc0d75457--`,
			statusExpected:   http.StatusOK,
			contentsExpected: "some data in a file",
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			ds.Reset()

			req, err := http.NewRequest("POST", "/",
				strings.NewReader(strings.ReplaceAll(tt.body, "\n", "\r\n")))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			if got, want := w.Code, tt.statusExpected; got != want {
				t.Fatalf("status=%d, want=%d", got, want)
			}
			if w.Code != http.StatusOK {
				return
			}

			for _, contents := range ds.entries {
				if got, want := contents, tt.contentsExpected; got != want {
					t.Fatalf("contents=%s, want=%s", got, want)
				}
			}
		})
	}
}
