package database

import "errors"

var (
	ErrAliasExists = errors.New("alias exists")
	ErrURLUnfound  = errors.New("url unfound")
	ErrIntOverflow = errors.New("integer overflow: aliasCount is out of range for uint32")
)
