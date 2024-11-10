package database

import (
	"context"
	"math"
	"strings"

	"github.com/5aradise/link-forge/internal/types"
	"github.com/5aradise/link-forge/internal/util"
)

func URLtoTypes(dbURL Url) types.URL {
	return types.URL{
		Id:    dbURL.ID,
		Alias: dbURL.Alias,
		Url:   dbURL.Url,
	}
}

type DB struct {
	q *Queries
}

func Create(db DBTX) *DB {
	return &DB{New(db)}
}

func (db *DB) CreateURL(ctx context.Context, alias, url string) (types.URL, error) {
	const op = "database.CreateURL"

	dbURL, err := db.q.CreateURL(ctx, CreateURLParams{
		Alias: alias,
		Url:   url,
	})
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return types.URL{}, util.OpWrap(op, ErrAliasExists)
		}

		return types.URL{}, util.OpWrap(op, err)
	}

	return URLtoTypes(dbURL), nil
}

func (db *DB) ListURLs(ctx context.Context) ([]types.URL, error) {
	const op = "database.ListURLs"

	dbURLs, err := db.q.ListURLs(ctx)
	if err != nil {
		return nil, util.OpWrap(op, err)
	}

	urls := make([]types.URL, 0, len(dbURLs))
	for _, dbURL := range dbURLs {
		urls = append(urls, URLtoTypes(dbURL))
	}
	return urls, nil
}

func (db *DB) GetURLByAlias(ctx context.Context, alias string) (types.URL, error) {
	const op = "database.GetURLByAlias"

	dbURL, err := db.q.GetURLByAlias(ctx, alias)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return types.URL{}, util.OpWrap(op, ErrURLUnfound)
		}

		return types.URL{}, util.OpWrap(op, err)
	}

	return URLtoTypes(dbURL), nil
}

func (db *DB) DeleteURLByAlias(ctx context.Context, alias string) (types.URL, error) {
	const op = "database.DeleteURLByAlias"

	dbURL, err := db.q.DeleteURLByAlias(ctx, alias)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return types.URL{}, util.OpWrap(op, ErrURLUnfound)
		}

		return types.URL{}, util.OpWrap(op, err)
	}

	return URLtoTypes(dbURL), nil
}

func (db *DB) LoadState(ctx context.Context) (uint32, error) {
	const op = "database.LoadState"

	aliasCount, err := db.q.LoadState(ctx)
	if err != nil {
		return 0, util.OpWrap(op, err)
	}
	if aliasCount < 0 || aliasCount > math.MaxUint32 {
		return 0, util.OpWrap(op, ErrIntOverflow)
	}

	return uint32(aliasCount), nil
}

func (db *DB) StoreState(ctx context.Context, aliasCount uint32) error {
	const op = "database.StoreState"

	err := db.q.StoreState(ctx, int64(aliasCount))
	if err != nil {
		return util.OpWrap(op, err)
	}
	return nil
}
