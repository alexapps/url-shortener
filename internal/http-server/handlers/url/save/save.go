package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "github.com/alexapps/url-shortener/internal/lib/api/response"
	sl "github.com/alexapps/url-shortener/internal/lib/logger/sl"
	"github.com/alexapps/url-shortener/internal/lib/random"
	"github.com/alexapps/url-shortener/internal/storage"
)

/*
	The save request handler
*/

// TODO: move to the config
const aliasLength = 6

// Struct describes the incommig save request
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// URLSaver is a interface that is implemented in the storage part
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(url string, alias string) (int64, error)
}

func New(l *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		l = l.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Incommig request targer struct object
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			l.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		l.Info("request body decoded", slog.Any("request", req))

		// Validate request using the 3rd party lib
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			l.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		// Check the alias
		alias := req.Alias
		if alias == "" {
			// the alias is not required parameter, if it ona missing, create one
			// TODO: check if generated alias is already exists
			alias = random.NewRandomString(aliasLength)
		}

		// save the incommig data
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			l.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url aleady exists"))

			return
		}
		if err != nil {
			l.Error("unable to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}
		l.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
