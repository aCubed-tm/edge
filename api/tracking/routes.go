package tracking

import (
	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/capture", addCapture)
	router.Get("/objects", getAllObjects)
	router.Get("/object/{uuid}", getObject)
	return router
}
