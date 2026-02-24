package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery оборачивает обработчик и перехватывает паники.
func Recovery(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered", "err", err, "stack", string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
