package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type mockStore struct {
	entries map[string]string
}

func (ds mockStore) GetEntry(id string) (string, error) {
	if contents, ok := ds.entries[id]; ok {
		return contents, nil
	}
	return "", errors.New("not found")
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
			bodyBytes, err := ioutil.ReadAll(w.Body)
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
			if tt.statusExpected != http.StatusOK {
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
	var pasteTests = []struct {
		contentType      string
		body             string
		statusExpected   int
		contentsExpected string
	}{
		{
			"text/plain",
			"hello, world!",
			http.StatusBadRequest,
			"",
		},
		// Empty input
		{
			"multipart/form-data; boundary=------------------------aea33768a2527972",
			`
--------------------------aea33768a2527972
Content-Disposition: form-data; name="dummyname1"



--------------------------aea33768a2527972--`,
			http.StatusBadRequest,
			"",
		},
		// Valid input
		{
			"multipart/form-data; boundary=------------------------aea33768a2527972",
			`
--------------------------aea33768a2527972
Content-Disposition: form-data; name="dummyname2"

some data I want to upload
--------------------------aea33768a2527972--`,
			http.StatusOK,
			"some data I want to upload",
		},
		{
			"multipart/form-data; boundary=------------------------ff01448fc0d75457",
			`
--------------------------ff01448fc0d75457
Content-Disposition: form-data; name="dummyname3"; filename="text.txt"
Content-Type: text/plain

some data in a file
--------------------------ff01448fc0d75457--`,
			http.StatusOK,
			"some data in a file",
		},
	}

	ds := mockStore{
		entries: make(map[string]string),
	}
	router := mux.NewRouter()
	s := defaultServer{
		store:  &ds,
		router: router,
	}
	s.routes()

	for _, tt := range pasteTests {
		ds.Reset()

		req, err := http.NewRequest("POST", "/",
			strings.NewReader(strings.ReplaceAll(tt.body, "\n", "\r\n")))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Content-Type", tt.contentType)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		if status := w.Code; status != tt.statusExpected {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, tt.statusExpected)
		}
		if tt.statusExpected != http.StatusOK {
			continue
		}
		for _, contents := range ds.entries {
			if contents != tt.contentsExpected {
				t.Fatalf("saved incorrect contents. got: [%s], want [%s]", contents, tt.contentsExpected)
			}
		}
	}
}
