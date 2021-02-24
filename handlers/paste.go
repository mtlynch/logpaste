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
		id := random.String(8)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read request body", http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err = s.store.InsertEntry(id, string(body))
		if err != nil {
			log.Printf("failed to save entry: %v", err)
			http.Error(w, "can't save entry", http.StatusInternalServerError)
			return
		}
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
