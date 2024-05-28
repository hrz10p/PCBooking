package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/pc_booking/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodGet, "/pc_booking/user/:id", app.getUserHandler)
	router.HandlerFunc(http.MethodGet, "/pc_booking/users/all", app.getAllUsersHandler)
	router.HandlerFunc(http.MethodDelete, "/pc_booking/user/:id/delete", app.requireAdminRole(app.deleteUserHandler))

	router.HandlerFunc(http.MethodPut, "/pc_booking/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/pc_booking/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic()
}
