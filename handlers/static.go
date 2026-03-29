package handlers

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

// CachingFileServer wraps an http.Handler to add caching headers
func CachingFileServer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add cache control headers
		w.Header().Set("Cache-Control", "public, max-age=604800") // 1 week

		// Let the original handler serve the file
		// The FileServer will handle Last-Modified and If-Modified-Since automatically
		h.ServeHTTP(w, r)
	})
}

// GetStaticFilesHandler returns an http.Handler that serves static files from the embedded filesystem
func GetStaticFilesHandler() http.Handler {
	// Get the static subdirectory as a filesystem
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	// Create a file server for the static files
	fileServer := http.FileServer(http.FS(staticFS))

	// Strip the /static/ prefix from the request path before passing to the file server
	// This is necessary because the URL path includes /static/ but our embedded filesystem
	// already has the files under the "static" directory
	handler := http.StripPrefix("/static/", fileServer)

	// Wrap the handler with our caching middleware
	return CachingFileServer(handler)
}
