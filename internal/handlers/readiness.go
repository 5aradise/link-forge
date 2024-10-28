package handlers

import (
	"log/slog"
	"net/http"

	"github.com/5aradise/link-forge/pkg/middleware"
)

func Readiness(l *slog.Logger) http.HandlerFunc {
	const op = "handlers.readiness"
	l = l.With(
		slog.String("op", op),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		l = l.With(
			slog.String("request_id", middleware.GetRequestID(r)),
		)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte{})
		if err != nil {
			l.Error("failed to write empty body", slog.String("error", err.Error()))
		}
	}
}
