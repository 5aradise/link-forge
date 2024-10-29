package database

import (
	"context"
	"fmt"
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
			return types.URL{}, fmt.Errorf("%s: %w", op, ErrAliasExists)
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
