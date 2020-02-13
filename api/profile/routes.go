package profile

import (
	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/user/{uuid}", getProfileUser)
	router.Put("/user/{uuid}", updateProfileUser)
	router.Post("/user/{uuid}", createProfileUser)

	return router
}
