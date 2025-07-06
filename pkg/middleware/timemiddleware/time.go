package timemiddleware

import (
	"net/http"
	"time"
)

type responseWriterWithHeader struct {
	http.ResponseWriter
	start       time.Time
	wroteHeader bool
}

func (rw *responseWriterWithHeader) WriteHeader(code int) {
	if !rw.wroteHeader {
		duration := time.Since(rw.start)
		rw.Header().Set("X-Response-Time", duration.String())
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWithHeader) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriterWithHeader{
			ResponseWriter: w,
			start:          time.Now(),
		}
		next.ServeHTTP(rw, r)
	})
}
