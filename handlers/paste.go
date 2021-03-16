package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mtlynch/logpaste/random"
)

const MaxPasteBytes = 2 * 1000 * 1000

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

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read request body", http.StatusBadRequest)
			return
		}

		if len(body) == 0 {
			log.Print("Paste body was empty")
			http.Error(w, "empty body", http.StatusBadRequest)
			return
		} else if len(body) > MaxPasteBytes {
			log.Printf("Paste body was too long: %d bytes", len(body))
			http.Error(w, "body too long", http.StatusBadRequest)
			return
		}

		id := random.String(8)
		err = s.store.InsertEntry(id, string(body))
		if err != nil {
			log.Printf("failed to save entry: %v", err)
			http.Error(w, "can't save entry", http.StatusInternalServerError)
			return
		}
		log.Printf("saved entry of %d bytes", len(body))

		w.Header().Set("Content-Type", "application/json")
		type response struct {
			Id string `json:"id"`
		}
		resp := response{
			Id: id,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}
