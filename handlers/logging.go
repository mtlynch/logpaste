package handlers

import (
	"io"
	"log"
	"net/http"
)

// LoggingHandler logs the method, path, status, and response size to out.
func LoggingHandler(out io.Writer, next http.Handler) http.Handler {
	logger := log.New(out, "", log.LstdFlags)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)
		logger.Printf("%s %q %d %d", r.Method, r.URL.Path, lw.status, lw.size)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	wroteHeader bool
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	if w.wroteHeader {
		return
	}
	// Informational responses do not commit the final status, except for a
	// protocol switch, which ends normal HTTP response handling.
	if code >= 100 && code < 200 && code != http.StatusSwitchingProtocols {
		return
	}
	w.status = code
	w.wroteHeader = true
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.status = http.StatusOK
		w.wroteHeader = true
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *loggingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *loggingResponseWriter) FlushError() error {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return http.NewResponseController(w.ResponseWriter).Flush()
}
