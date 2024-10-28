# link-forge

![Build Status](https://github.com/5aradise/learn-cicd-starter/actions/workflows/ci.yml/badge.svg)

## Description

RESTful API link shortener/customizer written in Go with ci/cd set up and using the best development approaches (code structuring, custom logging and middleware)

It uses [Standart http.ServeMux](https://pkg.go.dev/net/http@go1.23.2#ServeMux) as the HTTP router.

## Features
- Advanced custom logging
- Automated testing, style and security checks

## Technologies
- Go
- Docker

## Requirements
- Go 1.23.2+

## Local Development

Make sure you're on Go version 1.22+.

Create a copy of the `.env.example` file and rename it to `.env`.

Run the server:

```bash
make run
```