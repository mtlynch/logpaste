package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	muxHandlers "github.com/gorilla/handlers"

	"github.com/mtlynch/logpaste/handlers"
)

func main() {
	log.Print("Starting logpaste server")

	title := flag.String("title", "Log Paster", "title for the site")

	flag.Parse()

	s := handlers.New(handlers.SiteProperties{
		Title: *title,
	})
	http.Handle("/", muxHandlers.LoggingHandler(os.Stdout, s.Router()))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("Listening on %s", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
