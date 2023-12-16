package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	resp "github.com/alexapps/url-shortener/internal/lib/api/response"
	sl "github.com/alexapps/url-shortener/internal/lib/logger/sl"
	"github.com/alexapps/url-shortener/internal/storage"
)

type Requst struct {
}

type Response struct {
	resp.Response
	URL string `json:"url,omitempty"`
}

// URLGetter is an interface for getting url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(l *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		l = l.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			l.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))

		}

		// Getting URL value form storage
		u, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			l.Info("url not found", sl.Err(err))
			render.JSON(w, r, resp.Error("url not found"))

			return
		}
		if err != nil {
			l.Info("internal error", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		l.Info("got url", slog.String("url", u))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			URL:      u,
		})

		// or
		// redirect to found url
		// http.Redirect(w, r, resURL, http.StatusFound)
	}
}
