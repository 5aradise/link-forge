package handlers

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"log/slog"
	"net/http"

	"github.com/5aradise/link-forge/internal/database"
	"github.com/5aradise/link-forge/internal/types"
	"github.com/5aradise/link-forge/internal/util"
	"github.com/5aradise/link-forge/pkg/api"
	"github.com/5aradise/link-forge/pkg/middleware"
)

type URLStorage interface {
	CreateURL(ctx context.Context, alias, url string) (types.URL, error)
	ListURLs(ctx context.Context) ([]types.URL, error)
}

type URLService struct {
	l  *slog.Logger
	db URLStorage
}

func NewURLService(l *slog.Logger, db URLStorage) *URLService {
	return &URLService{
		l:  l,
		db: db,
	}
}

const (
	aliasSeq    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789()@:%_+.~#?&="
	aliasSeqLen = uint64(len(aliasSeq))
)

func (s *URLService) generateAlias() string {
	var (
		alias       string
		randomValue uint64
	)
	for range 6 {
		err := binary.Read(rand.Reader, binary.LittleEndian, &randomValue)
		if err != nil {
			return "abcdef"
		}
		alias += string(aliasSeq[randomValue%aliasSeqLen])
	}

	return alias
}

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
		writeJSONLog(w, http.StatusBadRequest, api.ResError(errMsg), l)
		return
	}

	if req.URL == "" {
		errMsg := "empty url field"
		l.Error("invalid request", slog.String("error", errMsg))
		writeJSONLog(w, http.StatusBadRequest, api.ResError(errMsg), l)
		return
	}

	l.Info("request body decoded", slog.Any("request", req))

	if !util.IsURL(req.URL) {
		errMsg := "invalid url"
		l.Error("invalid request", slog.String("error", errMsg), slog.String("url", req.URL))
		writeJSONLog(w, http.StatusBadRequest, api.ResError(errMsg), l)
		return
	}

	alias := req.Alias
	if alias == "" {
		alias = s.generateAlias()
	}

	newURL, err := s.db.CreateURL(r.Context(), alias, req.URL)
	if err != nil {
		if errors.Is(err, database.ErrAliasExists) {
			errMsg := "alias already exists"
			l.Info(errMsg, slog.String("alias", alias))
			writeJSONLog(w, http.StatusBadRequest, api.ResError(errMsg), l)
			return
		}

		errMsg := "failed to add url"
		l.Error(errMsg, util.SlErr(err))
		writeJSONLog(w, http.StatusInternalServerError, api.ResError(errMsg), l)
		return
	}

	l.Info("url added", slog.Int64("id", newURL.Id))

	writeJSONLog(w, http.StatusCreated, CreateURLResponse{
		api.ResOK(),
		newURL.Alias,
	}, l)
}

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
		writeJSONLog(w, http.StatusInternalServerError, api.ResError("failed to list urls"), l)
		return
	}

	l.Info("url listed")

	writeJSONLog(w, http.StatusCreated, ListURLsResponse{
		api.ResOK(),
		urls,
	}, l)
}
