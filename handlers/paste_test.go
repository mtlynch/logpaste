package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
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
	var pasteTests = []struct {
		id                 string
		httpStatusExpected int
		contentExpected    string
	}{
		// Valid entry
		{
			"12345678",
			http.StatusOK,
			"dummy entry",
		},
		// Non-existent entry
		{
			"missing1",
			http.StatusNotFound,
			"",
		},
	}

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

	for _, tt := range pasteTests {
		req, err := http.NewRequest("GET", "/"+tt.id, nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		if status := w.Code; status != tt.httpStatusExpected {
			t.Fatalf("for ID [%s], handler returned wrong status code: got %v want %v",
				tt.id, status, tt.httpStatusExpected)
		}
		if tt.httpStatusExpected != http.StatusOK {
			continue
		}
		bodyBytes, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Fatalf("failed to read HTTP response body: %v", err)
		}
		if tt.contentExpected != string(bodyBytes) {
			log.Fatalf("for ID [%s], got %s, want %s", tt.id, string(bodyBytes), tt.contentExpected)
		}
	}
}

func TestPastePut(t *testing.T) {
	var pasteTests = []struct {
		body               string
		httpStatusExpected int
	}{
		// Valid content
		{
			"hello, world!",
			http.StatusOK,
		},
		// Just at size limit
		{
			strings.Repeat("A", MaxPasteCharacters),
			http.StatusOK,
		},
		// Too long content
		{
			strings.Repeat("A", MaxPasteCharacters+1),
			http.StatusBadRequest,
		},
		// Empty content
		{
			"",
			http.StatusBadRequest,
		},
		// Just whitespace
		{
			"  ",
			http.StatusBadRequest,
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

		req, err := http.NewRequest("PUT", "/", strings.NewReader(tt.body))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		if status := w.Code; status != tt.httpStatusExpected {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, tt.httpStatusExpected)
		}
		if tt.httpStatusExpected != http.StatusOK {
			continue
		}
		resp := PastePutResponse{}
		err = json.NewDecoder(w.Body).Decode(&resp)
		if err != nil {
			t.Fatalf("failed to read HTTP response body: %v", err)
		}
		storedContents, err := ds.GetEntry(resp.Id)
		if err != nil {
			t.Fatalf("no entry found with id: %s", resp.Id)
		}
		if storedContents != tt.body {
			t.Fatalf("stored content doesn't match request body: stored: %s, want: %s", storedContents, tt.body)
		}
	}
}

func TestPastePost(t *testing.T) {
	var pasteTests = []struct {
		contentType        string
		body               string
		httpStatusExpected int
		contentsExpected   string
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
Content-Disposition: form-data; name="logpaste"



--------------------------aea33768a2527972--`,
			http.StatusBadRequest,
			"",
		},
		// Valid input
		{
			"multipart/form-data; boundary=------------------------aea33768a2527972",
			`
--------------------------aea33768a2527972
Content-Disposition: form-data; name="logpaste"

some data I want to upload
--------------------------aea33768a2527972--`,
			http.StatusOK,
			"some data I want to upload",
		},
		{
			"multipart/form-data; boundary=------------------------ff01448fc0d75457",
			`
--------------------------ff01448fc0d75457
Content-Disposition: form-data; name="logpaste"; filename="text.txt"
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

		if status := w.Code; status != tt.httpStatusExpected {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, tt.httpStatusExpected)
		}
		if tt.httpStatusExpected != http.StatusOK {
			continue
		}
		for _, contents := range ds.entries {
			if contents != tt.contentsExpected {
				t.Fatalf("saved incorrect contents. got: [%s], want [%s]", contents, tt.contentsExpected)
			}
		}
	}
}
