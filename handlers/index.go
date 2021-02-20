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
		type page struct {
			Title string
		}
		indexTemplate := template.Must(template.New(indexFilename).
			ParseFiles(path.Join(viewsRootDir, indexFilename)))
		if err := indexTemplate.ExecuteTemplate(w, indexFilename, page{
			Title: "LogPaste",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
