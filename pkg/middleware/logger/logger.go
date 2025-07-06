package logger

import (
	"errors"
	"go-monolite/pkg/logger"
	"go-monolite/pkg/middleware/request_id"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

func New(log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Str("request_id", request_id.GetReqID(r.Context())).
				Msg("incoming request")

			ww := NewWrapResponseWriter(w, r.ProtoMajor, 100)

			t1 := time.Now()
			defer func() {
				if rec := recover(); rec != nil {
					log.Error().
						Int("status", ww.Status()).
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Str("inf", ww.Body.String()).
						Str("request_id", request_id.GetReqID(r.Context())).
						Msg("log system error")

					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				var login *zerolog.Event
				status := ww.Status()
				// status == 400 - неправильный ввод данных
				if status >= 200 && status < 300 {
					login = log.Info()
				} else {
					login = log.Warn().Err(errors.New("request failed"))
				}
				login.
					Int("status", ww.Status()).
					Str("inf", ww.Body.String()).
					Str("elapsed_ms", time.Since(t1).String()).
					Str("request_id", request_id.GetReqID(r.Context())).
					Msg("request processed")
			}()

			// добавление Logger в Context
			ctx := logger.WithContext(r.Context(), "request_id", request_id.GetReqID(r.Context()))

			next.ServeHTTP(ww, r.WithContext(ctx))

		})

		return http.HandlerFunc(fn)
	}
}
