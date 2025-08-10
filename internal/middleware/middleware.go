package middleware

import (
	"net/http"
	"runtime/debug"
	"strings"

	"go.uber.org/zap"
	
	"github.com/sunr3d/subscription-aggregator/internal/httpx"
)

func JSONValidator(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch:
				ct := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
				if !httpx.IsJSON(ct) {
					if err := httpx.HttpError(w, http.StatusUnsupportedMediaType, "Ожидается Content-Type: application/json"); err != nil {
						log.Warn("JSONValidator: не удалось записать ошибку",
							zap.Error(err),
							zap.String("method", r.Method),
							zap.String("url", r.URL.Path),
							zap.String("content_type", ct),
						)
					}
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Recovery(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("Паника в обработчике запроса",
						zap.Any("rec", rec),
						zap.String("stack", string(debug.Stack())),
						zap.String("url", r.URL.Path),
						zap.String("method", r.Method),
					)
					if err := httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера"); err != nil {
						log.Warn("Recovery: не удалось записать ошибку",
							zap.Error(err),
							zap.String("method", r.Method),
							zap.String("url", r.URL.Path),
						)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
