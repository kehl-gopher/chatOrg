package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routers() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/api/v1/company", app.RegisterCompany)
	router.HandlerFunc(http.MethodPost, "/api/v1/about", app.AboutEndpoint)
	router.HandlerFunc(http.MethodPost, "/api/v1/chat", app.HandleChat)

	// router.HandlerFunc(http.MethodPost, "/api/v1/webhook", app.webHook)

	router.HandlerFunc(http.MethodGet, "/integrations", app.appIntegration)
	router.HandlerFunc(http.MethodGet, "/ping", ping)
	return app.recoverMiddleware(corsMiddleWare(router))
}
