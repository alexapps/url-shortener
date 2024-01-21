package remove

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	resp "github.com/alexapps/url-shortener/internal/lib/api/response"
	sl "github.com/alexapps/url-shortener/internal/lib/logger/sl"
)

// Struct describes the incommig save request
type Request struct {
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLRemover
type URLRemover interface {
	DeleteURL(alias string) error
}

func New(l *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.remove.New"

		l = l.With(
			slog.String("op", op),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, req)
		if err != nil {
			l.Error("failed to decode request body", sl.Err(err))
		}

		l.Info("request body decoded", slog.Any("request", req))

		err = urlRemover.DeleteURL(req.Alias)
		if err != nil {
			l.Error("unable to delete url by alias", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete url"))

			return
		}
		l.Info("url deleted", slog.String("alias", req.Alias))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
