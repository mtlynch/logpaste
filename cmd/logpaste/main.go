package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	gorilla "github.com/mtlynch/gorilla-handlers"

	"github.com/mtlynch/logpaste/handlers"
)

func main() {
	log.Print("Starting logpaste server")

	perMinuteLimit := flag.Int("perminutelimit",
		0, "number of pastes to allow per IP per minute (set to 0 to disable rate limiting)")
	maxPasteMiB := flag.Int64("maxsize", 2, "max file size as MiB")

	flag.Parse()

	const charactersPerMiB = 1024 * 1024
	maxCharLimit := *maxPasteMiB * charactersPerMiB

	h := gorilla.LoggingHandler(os.Stdout, handlers.New(*perMinuteLimit, maxCharLimit).Router())
	if os.Getenv("LP_BEHIND_PROXY") != "" {
		h = gorilla.ProxyIPHeadersHandler(h)
	}
	http.Handle("/", h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("Listening on %s", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
