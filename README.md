# link-forge

![Build Status](https://github.com/5aradise/link-forge/actions/workflows/ci.yml/badge.svg)

## Description

RESTful API link shortener/customizer written in Go with ci/cd set up and using the best development approaches (code structuring, advanced testing with mocking, custom logging and middleware)

It uses [Standart http.ServeMux](https://pkg.go.dev/net/http@go1.23.2#ServeMux) as the HTTP router and [Turso](https://turso.tech)(SQLite) as the cloud database with [goose](https://github.com/pressly/goose) for migration and [sqlc](https://sqlc.dev/) for compile SQL queries.

## Features

- Advanced custom logging
- Automated testing with mocking, style and security checks

## Technologies

- Go
- Testify
- Mockery
- Goose
- Sqlc
- SQLite(Turso)
- Docker

## Requirements

- Go 1.23.2+

## Local Development

Make sure you're on Go version 1.23.2+

Create a copy of the `.env.example` file and rename it to `.env`

### In `cmd/link-forge.go`:

Recomment import lines:

```go
// _ "github.com/tursodatabase/libsql-client-go/libsql"
_ "github.com/mattn/go-sqlite3"
```

And update sql.Open line:

```go
conn, err := sql.Open("libsql", config.Cfg.DB.URL)   // ->
conn, err := sql.Open("sqlite3", config.Cfg.DB.URL)  // <-
```

### Install dependencies:

```bash
go mod download
```

### In `scripts/`:

Recomment goose lines:

```bash
# goose turso $DATABASE_URL ...
goose sqlite3 $DATABASE_URL ...
```

### Install goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Run migrations:

```bash
./scripts/migrateup.sh
```

### Run the server:

```bash
make run
```

or

```bash
go build -C cmd/link-forge/ -o ../../bin/link-forge && CONFIG_PATH=.env ./bin/link-forge
```
