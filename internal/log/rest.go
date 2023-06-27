package log

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func RestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := CreateContext(r.Context(), "request_id", middleware.GetReqID(r.Context()))
		r = r.WithContext(ctx)

		ts := time.Now()
		WithContext(ctx).Trace().Msgf("Entry --> : %s %s", r.Method, r.RequestURI)

		myWritter := &customWriter{w, 200}
		next.ServeHTTP(myWritter, r)

		WithContext(ctx).Trace().
			Str("latency", time.Since(ts).String()).
			Int("status", myWritter.statusCode).
			Msgf("Exit <-- : %s %s", r.Method, r.RequestURI)
	})
}

type customWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *customWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *customWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
