package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestLoggingResponseWriterBookkeeping(t *testing.T) {
	for _, tt := range []struct {
		name   string
		write  func(*loggingResponseWriter)
		status int
		size   int
	}{
		{
			name:   "empty response defaults to OK",
			write:  func(w *loggingResponseWriter) {},
			status: http.StatusOK,
		},
		{
			name: "implicit body",
			write: func(w *loggingResponseWriter) {
				_, _ = w.Write([]byte("hello"))
			},
			status: http.StatusOK,
			size:   5,
		},
		{
			name: "explicit error body",
			write: func(w *loggingResponseWriter) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("bad"))
			},
			status: http.StatusBadRequest,
			size:   3,
		},
		{
			name: "no body",
			write: func(w *loggingResponseWriter) {
				w.WriteHeader(http.StatusNoContent)
			},
			status: http.StatusNoContent,
		},
		{
			name: "repeated final header keeps first",
			write: func(w *loggingResponseWriter) {
				w.WriteHeader(http.StatusCreated)
				w.WriteHeader(http.StatusBadRequest)
			},
			status: http.StatusCreated,
		},
		{
			name: "header after body ignored",
			write: func(w *loggingResponseWriter) {
				_, _ = w.Write([]byte("body"))
				w.WriteHeader(http.StatusInternalServerError)
			},
			status: http.StatusOK,
			size:   4,
		},
		{
			name: "informational then final",
			write: func(w *loggingResponseWriter) {
				w.WriteHeader(http.StatusEarlyHints)
				w.WriteHeader(http.StatusCreated)
			},
			status: http.StatusCreated,
		},
		{
			name: "protocol switch is final",
			write: func(w *loggingResponseWriter) {
				w.WriteHeader(http.StatusSwitchingProtocols)
				w.WriteHeader(http.StatusOK)
			},
			status: http.StatusSwitchingProtocols,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			underlying := httptest.NewRecorder()
			rec := &loggingResponseWriter{ResponseWriter: underlying, status: http.StatusOK}
			tt.write(rec)
			if got, want := rec.status, tt.status; got != want {
				t.Errorf("status=%d, want=%d", got, want)
			}
			if got, want := rec.size, tt.size; got != want {
				t.Errorf("size=%d, want=%d", got, want)
			}
		})
	}
}

type partialErrorWriter struct {
	max int
	err error
}

func (w partialErrorWriter) Header() http.Header         { return make(http.Header) }
func (w partialErrorWriter) WriteHeader(int)             {}
func (w partialErrorWriter) Write(p []byte) (int, error) { return w.max, w.err }

func TestLoggingResponseWriterCountsPartialWrite(t *testing.T) {
	writeError := errors.New("write failed")
	rec := &loggingResponseWriter{ResponseWriter: partialErrorWriter{max: 2, err: writeError}}
	n, err := rec.Write([]byte("hello"))
	if got, want := n, 2; got != want {
		t.Errorf("Write() bytes=%d, want=%d", got, want)
	}
	if !errors.Is(err, writeError) {
		t.Fatalf("Write() error=%v, want=%v", err, writeError)
	}
	if got, want := rec.size, 2; got != want {
		t.Errorf("size=%d, want=%d", got, want)
	}
}

func TestLoggingResponseWriterContentTypeSniffing(t *testing.T) {
	underlying := httptest.NewRecorder()
	rec := &loggingResponseWriter{ResponseWriter: underlying}
	_, _ = rec.Write([]byte("hello"))
	if got, want := underlying.Header().Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("Content-Type=%q, want=%q", got, want)
	}
}

func TestLoggingHandlerForwardsResponse(t *testing.T) {
	h := LoggingHandler(io.Discard, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "yes")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, "created")
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if got, want := rec.Code, http.StatusCreated; got != want {
		t.Fatalf("status=%d, want=%d", got, want)
	}
	if got, want := rec.Body.String(), "created"; got != want {
		t.Errorf("body=%q, want=%q", got, want)
	}
	if got, want := rec.Header().Get("X-Test"), "yes"; got != want {
		t.Errorf("X-Test=%q, want=%q", got, want)
	}
}

func TestLoggingResponseWriterFlush(t *testing.T) {
	underlying := httptest.NewRecorder()
	rec := &loggingResponseWriter{ResponseWriter: underlying}
	if err := http.NewResponseController(rec).Flush(); err != nil {
		t.Fatalf("Flush()=%v, want nil", err)
	}
	rec.WriteHeader(http.StatusInternalServerError)
	if got, want := rec.status, http.StatusOK; got != want {
		t.Errorf("status after flush and late header=%d, want=%d", got, want)
	}
	if got, want := underlying.Code, http.StatusOK; got != want {
		t.Errorf("underlying status=%d, want=%d", got, want)
	}
	if !underlying.Flushed {
		t.Error("underlying writer was not flushed")
	}
	if got, want := rec.Unwrap(), underlying; got != want {
		t.Error("Unwrap() did not return the underlying writer")
	}
}

type headerRecordingWriter struct {
	*httptest.ResponseRecorder
	codes []int
}

func (w *headerRecordingWriter) WriteHeader(code int) {
	w.codes = append(w.codes, code)
	w.ResponseRecorder.WriteHeader(code)
}

func TestLoggingResponseWriterForwardsRepeatedHeaders(t *testing.T) {
	underlying := &headerRecordingWriter{ResponseRecorder: httptest.NewRecorder()}
	rec := &loggingResponseWriter{ResponseWriter: underlying}
	rec.WriteHeader(http.StatusCreated)
	rec.WriteHeader(http.StatusInternalServerError)
	if got, want := underlying.codes, []int{http.StatusCreated, http.StatusInternalServerError}; !slices.Equal(got, want) {
		t.Errorf("forwarded status codes=%v, want=%v", got, want)
	}
	if got, want := rec.status, http.StatusCreated; got != want {
		t.Errorf("recorded status=%d, want=%d", got, want)
	}
	if got, want := underlying.Code, http.StatusCreated; got != want {
		t.Errorf("underlying status=%d, want=%d", got, want)
	}
}
