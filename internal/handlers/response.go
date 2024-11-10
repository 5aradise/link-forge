package handlers

import (
	"log/slog"
	"net/http"

	"github.com/5aradise/link-forge/internal/util"
	"github.com/5aradise/link-forge/pkg/api"
)

func WriteJSONLog(w http.ResponseWriter, statusCode int, v any, l *slog.Logger) {
	err := api.WriteJSON(w, statusCode, v)
	if err != nil {
		l.Error("failed to write response", util.SlErr(err))
	}
}
