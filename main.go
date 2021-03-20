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

	title := flag.String("title", "LogPaste", "title for the site")
	subtitle := flag.String("subtitle",
		"A minimalist, open-source debug log upload service",
		"subtitle for the site")
	showDocs := flag.Bool("showdocs",
		true, "whether to display usage information on homepage")

	flag.Parse()

	s := handlers.New(handlers.SiteProperties{
		Title:    *title,
		Subtitle: *subtitle,
		ShowDocs: *showDocs,
	})
	http.Handle("/", muxHandlers.LoggingHandler(os.Stdout, s.Router()))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("Listening on %s", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
