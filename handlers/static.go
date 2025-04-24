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

	// Wrap the file server with our caching middleware
	return CachingFileServer(http.FileServer(http.FS(staticFS)))
}
