package remove

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
)

type URLRemover interface {
	RemoveURL(alias string) error
}

func New(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.remove"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlRemover.RemoveURL(alias)
		if err != nil {
			log.Error("failed to delete", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("removed url", slog.String("alias", alias))

		render.JSON(w, r, resp.OK())
	}
}
