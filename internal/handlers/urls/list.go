package urls

import (
	"log/slog"
	"net/http"

	"github.com/5aradise/link-forge/internal/handlers"
	"github.com/5aradise/link-forge/internal/types"
	"github.com/5aradise/link-forge/internal/util"
	"github.com/5aradise/link-forge/pkg/api"
	"github.com/5aradise/link-forge/pkg/middleware"
)

type ListURLsResponse struct {
	api.Response
	URLs []types.URL `json:"urls"`
}

func (s *URLService) ListURLs(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.list"

	l := s.l.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestID(r)),
	)

	urls, err := s.db.ListURLs(r.Context())
	if err != nil {
		l.Error("failed to list urls", util.SlErr(err))
		handlers.WriteJSONLog(w, http.StatusInternalServerError, api.ResError("failed to list urls"), l)
		return
	}

	l.Info("urls listed")

	handlers.WriteJSONLog(w, http.StatusOK, ListURLsResponse{
		api.ResOK(),
		urls,
	}, l)
}
