package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	appLogger "store-service/internal/logger"
)

// LoggerMiddleware attaches zap logger with request metadata into context and logs request summary.
func LoggerMiddleware(base *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			log := base.With(
				zap.String("req_id", reqID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)
			ctx := appLogger.WithContext(r.Context(), log)
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r.WithContext(ctx))

			log.Info("request completed",
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
