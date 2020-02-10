package auth

import (
	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/authenticate", authenticate)
	router.Post("/register", register)
	router.Post("/meet", getUserUuidAndInvites) // used to be at /check-registration
	return router
}
