package profile

import (
	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/user/{uuid}", getProfileUser)
	router.Put("/user/{uuid}", updateProfileUser)
	router.Post("/user/{uuid}", createProfileUser)
	router.Get("/organisation/{uuid}", getProfileOrganization)
	router.Put("/organisation/{uuid}", updateProfileOrganization)
	router.Post("/organisation/{uuid}", createProfileOrganization)

	router.Get("/user/{uuid}/emails", getUserEmails)

	return router
}
