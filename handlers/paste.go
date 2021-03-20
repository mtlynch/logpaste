package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/mtlynch/logpaste/random"
)

const MaxPasteCharacters = 2 * 1000 * 1000

type PastePutResponse struct {
	Id string `json:"id"`
}

func (s defaultServer) pasteGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		contents, err := s.store.GetEntry(id)
		if err != nil {
			log.Printf("Error retrieving entry with id %s: %v", id, err)
			http.Error(w, "entry not found", http.StatusNotFound)
		}
		io.WriteString(w, contents)
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

		bodyRaw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read request body", http.StatusBadRequest)
			return
		}

		body := string(bodyRaw)
		if !validatePaste(body, w) {
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
		err := r.ParseMultipartForm(MaxPasteCharacters + 1024)
		if err != nil {
			log.Printf("failed to parse form: %v", err)
			http.Error(w, "no valid multipart/form-data found", http.StatusBadRequest)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		formValues, ok := r.MultipartForm.Value["logpaste"]
		if !ok {
			log.Print("Form did not contain expected field: logpaste")
			http.Error(w, "logpaste form data is required", http.StatusBadRequest)
		}
		if len(formValues) < 1 {
			log.Print("logpaste form data contains no values")
			http.Error(w, "logpaste form data in unexpected format", http.StatusBadRequest)
		}

		body := formValues[0]
		if !validatePaste(body, w) {
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

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("http://%s/%s\n", r.Host, id)))
	}
}

func generateEntryId() string {
	return random.String(8)
}

func validatePaste(p string, w http.ResponseWriter) bool {
	if len(strings.TrimSpace(p)) == 0 {
		log.Print("Paste body was empty")
		http.Error(w, "empty body", http.StatusBadRequest)
		return false
	} else if len(p) > MaxPasteCharacters {
		log.Printf("Paste body was too long: %d characters", len(p))
		http.Error(w, "body too long", http.StatusBadRequest)
		return false
	}
	return true
}
