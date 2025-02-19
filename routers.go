package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func routers() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/api/v1/about", createFaq)

	router.HandlerFunc(http.MethodPost, "/api/v1/webhook", webHook)

	router.HandlerFunc(http.MethodGet, "/integrations", appIntegration)
	return router
}
