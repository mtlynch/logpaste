package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"github.com/mtlynch/logpaste/random"
	"github.com/mtlynch/logpaste/store"
)

type PastePutResponse struct {
	Id string `json:"id"`
}

func (s defaultServer) pasteGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		contents, err := s.store.GetEntry(id)
		if err != nil {
			if _, ok := err.(store.EntryNotFoundError); ok {
				http.Error(w, "entry not found", http.StatusNotFound)
				return
			}
			log.Printf("failed to retrieve entry %s from datastore: %v", id, err)
			http.Error(w, "failed to retrieve entry", http.StatusInternalServerError)
			return
		}
		if _, err := io.WriteString(w, contents); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	}
}

func (s defaultServer) pasteOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}

func (s defaultServer) pastePut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		bodyRaw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, s.maxCharLimit))
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read request body", http.StatusBadRequest)
			return
		}

		body := string(bodyRaw)
		if !validatePaste(body, w, s.maxCharLimit) {
			return
		}

		id := generateEntryId()
		err = s.store.InsertEntry(id, body)
		if err != nil {
			log.Printf("failed to save entry: %v", err)
			http.Error(w, "can't save entry", http.StatusInternalServerError)
			return
		}
		log.Printf("saved entry of %d characters", len(body))

		w.Header().Set("Content-Type", "application/json")
		resp := PastePutResponse{
			Id: id,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}

func (s defaultServer) pastePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, s.maxCharLimit+1024)
		if err := r.ParseMultipartForm(s.maxCharLimit); err != nil {
			log.Printf("failed to parse form: %v", err)
			http.Error(w, "no valid multipart/form-data found", http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")

		body, ok := parsePasteFromMultipartForm(r.MultipartForm, w)
		if !ok {
			log.Print("form did not contain any recognizable data")
			http.Error(w, "form data or file is required", http.StatusBadRequest)
			return
		}

		if !validatePaste(body, w, s.maxCharLimit) {
			return
		}

		id := generateEntryId()
		if err := s.store.InsertEntry(id, body); err != nil {
			log.Printf("failed to save entry: %v", err)
			http.Error(w, "can't save entry", http.StatusInternalServerError)
			return
		}
		log.Printf("saved entry of %d characters", len(body))

		w.Header().Set("Content-Type", "text/plain")
		resultURL := fmt.Sprintf("%s/%s\n", baseURLFromRequest(r), id)
		if _, err := w.Write([]byte(resultURL)); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	}
}

func generateEntryId() string {
	return random.String(8)
}

func validatePaste(p string, w http.ResponseWriter, maxCharLimit int64) bool {
	if len(strings.TrimSpace(p)) == 0 {
		log.Print("Paste body was empty")
		http.Error(w, "empty body", http.StatusBadRequest)
		return false
	} else if int64(len(p)) > maxCharLimit {
		log.Printf("Paste body was too long: %d characters", len(p))
		http.Error(w, "body too long", http.StatusBadRequest)
		return false
	}
	return true
}

func parsePasteFromMultipartForm(f *multipart.Form, w http.ResponseWriter) (string, bool) {
	if content, ok := parsePasteFromMultipartFormValue(f); ok {
		return content, true
	}
	if content, ok := parsePasteFromMultipartFormFile(f, w); ok {
		return content, true
	}
	return "", false
}

func parsePasteFromMultipartFormValue(f *multipart.Form) (string, bool) {
	return anyValueInForm(f)
}

func anyValueInForm(f *multipart.Form) (string, bool) {
	for _, values := range f.Value {
		if len(values) < 1 {
			log.Printf("form values are empty")
			continue
		}
		return values[0], true
	}
	return "", false
}

func parsePasteFromMultipartFormFile(f *multipart.Form, w http.ResponseWriter) (string, bool) {
	file, ok := anyFileInForm(f.File)
	if !ok {
		return "", false
	}

	body, err := io.ReadAll(file)
	if err != nil {
		log.Printf("failed to read form file: %v", err)
		return "", false
	}

	return string(body), true
}

func anyFileInForm(formFiles map[string][]*multipart.FileHeader) (multipart.File, bool) {
	for _, fileHeaders := range formFiles {
		if len(fileHeaders) < 1 {
			log.Printf("form files are empty")
			continue
		}
		file, err := fileHeaders[0].Open()
		if err != nil {
			log.Printf("failed to open form file: %v", err)
			return nil, false
		}
		return file, true
	}
	return nil, false
}

func baseURLFromRequest(r *http.Request) string {
	var scheme string
	// If we're running behind a proxy, assume that it's a TLS proxy.
	if r.TLS != nil || os.Getenv("LP_BEHIND_PROXY") != "" {
		scheme = "https"
	} else {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}
