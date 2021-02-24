package handlers

import (
	"log"
	"net/http"
	"os"
	"path"
)

const staticRootDir = "./static"

func (s defaultServer) serveStaticResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fs := http.Dir(staticRootDir)
		file, err := fs.Open(r.URL.Path)
		if os.IsNotExist(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Failed to retrieve the file %s from the file system: %s", r.URL.Path, err)
			http.Error(w, "Failed to find file: "+r.URL.Path, http.StatusInternalServerError)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			log.Printf("Failed to retrieve the information of %s from the file system: %s", r.URL.Path, err)
			http.Error(w, "Failed to serve: "+r.URL.Path, http.StatusInternalServerError)
			return
		}
		if stat.IsDir() {
			log.Printf("client attempted to access a static directory: %s", r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		// Otherwise, serve a static file.
		http.ServeFile(w, r, path.Join(staticRootDir, r.URL.Path))
	}
}
