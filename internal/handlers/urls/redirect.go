package urls

import (
	"log/slog"
	"net/http"

	"github.com/5aradise/link-forge/internal/util"
	"github.com/5aradise/link-forge/pkg/api"
	"github.com/5aradise/link-forge/pkg/middleware"
)

func (s *URLService) RedirectURL(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.redirect"

	const PageNotFoundHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Page Not Found | 404</title>
  </head>
  <body>
    Bad link
  </body>
</html>`

	l := s.l.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestID(r)),
	)

	alias := r.PathValue("alias")
	if alias == "" {
		panic("empty alias path value")
	}

	url, err := s.db.GetURLByAlias(r.Context(), alias)
	if err != nil {
		l.Error("failed to get url", util.SlErr(err))
		err := api.WriteHTML(w, http.StatusNotFound, PageNotFoundHTML)
		if err != nil {
			l.Error("failed to write response", util.SlErr(err))
		}
		return
	}

	l.Info("redirected to url", slog.String("url", url.Url))

	http.Redirect(w, r, url.Url, http.StatusFound)
}
