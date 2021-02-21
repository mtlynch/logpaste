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
			BackendBaseURL string
		}
		indexTemplate := template.Must(template.New(indexFilename).
			ParseFiles(path.Join(viewsRootDir, indexFilename)))
		if err := indexTemplate.ExecuteTemplate(w, indexFilename, page{
			//BackendBaseURL: "https://logs.tinypilotkvm.com",
			BackendBaseURL: "https://logpaste-fqju5wfryq-ue.a.run.app",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
