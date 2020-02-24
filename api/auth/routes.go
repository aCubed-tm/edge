package auth

import (
	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/authenticate", authenticate)
	router.Post("/register", register)
	router.Post("/meet", getUserUuidAndInvites) // used to be at /check-registration
	router.Get("/activate/{token}", verifyEmail)
	router.Get("/close", dropCurrentToken)
	router.Get("/logout", dropAllTokens)

	router.Put("/email/{uuid}", putEmail)
	router.Post("/email/", addEmail)
	router.Delete("/email/{uuid}", deleteEmail)

	return router
}
