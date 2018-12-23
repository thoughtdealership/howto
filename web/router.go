package web

import (
	"net/http"

	"github.com/thoughtdealership/howto/web/middle"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Create generates the http router
func Create() http.Handler {
	router := httprouter.New()

	router.Handle("GET", "/hello", middle.ResponseHandler(Hello))
	router.Handle("GET", "/error", middle.ResponseHandler(ErrorRoute))
	router.Handle("GET", "/user-error", middle.ResponseHandler(UserErrorRoute))
	router.Handle("GET", "/multi-error", middle.ResponseHandler(MultiErrorRoute))
	router.Handle("GET", "/panic", middle.ResponseHandler(Panic))
	router.Handle("GET", "/version", Version)

	return alice.New(
		middle.RecoveryHandler,
		middle.FrameHandler,
		middle.RequestIDHandler,
		middle.RequestPathHandler,
		middle.BodyHandler).
		Then(router)
}
