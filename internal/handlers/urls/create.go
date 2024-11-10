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

type CreateURLRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type CreateURLResponse struct {
	api.Response
	Alias string `json:"alias,omitempty"`
}

func (s *URLService) CreateURL(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.create"
	l := s.l.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestID(r)),
	)

	var req CreateURLRequest
	if err := api.DecodeJSON(r, &req); err != nil {
		errMsg := "failed to decode request body"
		l.Error(errMsg, util.SlErr(err))
		handlers.WriteJSONLog(w, http.StatusBadRequest, api.ResError(errMsg), l)
		return
	}

	if req.URL == "" {
		errMsg := "empty url field"
		l.Error("invalid request", slog.String("error", errMsg))
		handlers.WriteJSONLog(w, http.StatusBadRequest, api.ResError(errMsg), l)
		return
	}

	l.Info("request body decoded", slog.Any("request", req))

	if !util.IsURL(req.URL) {
		errMsg := "invalid url"
		l.Error("invalid request", slog.String("error", errMsg), slog.String("url", req.URL))
		handlers.WriteJSONLog(w, http.StatusBadRequest, api.ResError(errMsg), l)
		return
	}

	alias := req.Alias
	if alias == "" {
		var err error
		alias, err = s.as.nextAlias()
		if err != nil {
			l.Error("ALIAS COUNT IS EXCEEDED")
			handlers.WriteJSONLog(w, http.StatusInternalServerError, api.ResError("failed to generate alias"), l)
			return
		}

		l.Info("generated new alias", slog.String("alias", alias))
	} else if len(alias) <= maxAliasLen {
		errMsg := "alias length is too short"
		l.Info(errMsg, slog.String("alias", alias))
		handlers.WriteJSONLog(w, http.StatusBadRequest, api.ResError(errMsg), l)
		return
	}

	newURL, err := s.db.CreateURL(r.Context(), alias, req.URL)
	if err != nil {
		if errors.Is(err, database.ErrAliasExists) {
			errMsg := "alias already exists"
			l.Info(errMsg, slog.String("alias", alias))
			handlers.WriteJSONLog(w, http.StatusBadRequest, api.ResError(errMsg), l)
			return
		}

		errMsg := "failed to add url"
		l.Error(errMsg, util.SlErr(err))
		handlers.WriteJSONLog(w, http.StatusInternalServerError, api.ResError(errMsg), l)
		return
	}

	l.Info("url added", slog.Int64("id", newURL.Id))

	handlers.WriteJSONLog(w, http.StatusCreated, CreateURLResponse{
		api.ResOK(),
		newURL.Alias,
	}, l)
}
