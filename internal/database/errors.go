package database

import "errors"

var (
	ErrAliasExists = errors.New("alias exists")
	ErrURLUnfound  = errors.New("url unfound")
)
