package urls

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/5aradise/link-forge/internal/database"
	"github.com/5aradise/link-forge/internal/handlers"
	"github.com/5aradise/link-forge/internal/util"
	"github.com/5aradise/link-forge/pkg/api"
	"github.com/5aradise/link-forge/pkg/middleware"
)

func (s *URLService) DeleteURL(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.redirect"

	l := s.l.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestID(r)),
	)

	alias := r.PathValue("alias")
	if alias == "" {
		panic("empty alias path value")
	}

	url, err := s.db.DeleteURLByAlias(r.Context(), alias)
	if err != nil {
		l.Error("failed to delete url", util.SlErr(err))
		if errors.Is(err, database.ErrURLUnfound) {
			handlers.WriteJSONLog(w, http.StatusBadRequest, api.ResError("url with this alias unfound"), l)
		} else {
			handlers.WriteJSONLog(w, http.StatusInternalServerError, api.ResError("internal error"), l)
		}
		return
	}

	l.Info("url deleted", slog.Any("url", url))

	handlers.WriteJSONLog(w, http.StatusOK, api.ResOK(), l)
}
