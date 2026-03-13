package middleware

import (
	"net/http"
	"time"

	"socialai/logger"
	"go.uber.org/zap"
	"github.com/gorilla/mux"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		logger.Logger.Info("request",
			zap.String("method", r.Method),
			zap.String("path", path),
			zap.Int("status", recorder.status),
			zap.Float64("latency_ms",  float64(time.Since(startTime).Milliseconds())),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}