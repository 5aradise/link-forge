package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type LogInformer interface {
	Info(msg string, args ...any)
}

type LogErrorer interface {
	Error(msg string, args ...any)
}
