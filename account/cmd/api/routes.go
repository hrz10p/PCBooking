package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/pc_booking/user/register/", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/pc_booking/user/login/", app.loginUserHandler)

	router.HandlerFunc(http.MethodGet, "/pc_booking/user/", app.getUserByEmailHandler)
	router.HandlerFunc(http.MethodGet, "/pc_booking/users/all/", app.getAllUsersHandler)
	router.HandlerFunc(http.MethodDelete, "/pc_booking/user/delete/", app.requireAdminRole(app.deleteByEmailUserHandler))

	router.HandlerFunc(http.MethodPut, "/pc_booking/users/activated/", app.activateUserHandler)

	return app.recoverPanic(router)
}
