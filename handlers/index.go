package handlers

import (
	"net/http"
	"path"
	"text/template"
)

const indexFilename = "index.html"
const viewsRootDir = "./views"

func (s defaultServer) serveIndexPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexTemplate := template.Must(template.New(indexFilename).
			ParseFiles(path.Join(viewsRootDir, indexFilename)))
		if err := indexTemplate.ExecuteTemplate(w, indexFilename, s.siteProps); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
