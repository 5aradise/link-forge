package urls

import (
	"context"
	"log/slog"

	"github.com/5aradise/link-forge/internal/types"
	"github.com/5aradise/link-forge/internal/util"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --name=URLStorage
type URLStorage interface {
	CreateURL(ctx context.Context, alias, url string) (types.URL, error)
	ListURLs(ctx context.Context) ([]types.URL, error)
	GetURLByAlias(ctx context.Context, alias string) (types.URL, error)
	DeleteURLByAlias(ctx context.Context, alias string) (types.URL, error)
}

type URLService struct {
	l  *slog.Logger
	db URLStorage
	as aliasService
}

func NewService(l *slog.Logger, db URLStorage, currAliasCount uint32) (*URLService, error) {
	const op = "handlers.url.NewService"
	as, err := newAliasService(currAliasCount)
	if err != nil {
		return nil, util.OpWrap(op, err)
	}

	return &URLService{
		l:  l,
		db: db,
		as: as,
	}, nil
}

func (s *URLService) AliasCount() uint32 {
	return s.as.loadCount()
}
