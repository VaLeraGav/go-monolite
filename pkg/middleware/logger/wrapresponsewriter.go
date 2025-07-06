package logger

import (
	"bytes"
	"net/http"
)

type WrapResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Body       *bytes.Buffer
	MaxSize    int
}

func NewWrapResponseWriter(w http.ResponseWriter, protoMajor int, maxSize int) *WrapResponseWriter {
	return &WrapResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		Body:           new(bytes.Buffer),
		MaxSize:        maxSize,
	}
}

func (ww *WrapResponseWriter) Write(b []byte) (int, error) {
	n, err := ww.ResponseWriter.Write(b)

	if ww.Body.Len() < ww.MaxSize {
		remaining := ww.MaxSize - ww.Body.Len()

		if len(b) <= remaining {
			ww.Body.Write(b)
		} else {
			allowedLen := remaining - 3
			if allowedLen > 0 {
				ww.Body.Write(b[:allowedLen])
			}
			ww.Body.WriteString("...")
		}
	}

	return n, err
}

func (ww *WrapResponseWriter) WriteHeader(statusCode int) {
	ww.StatusCode = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}

func (ww *WrapResponseWriter) Status() int {
	return ww.StatusCode
}
