package routes

import (
	"proxyScanner/Server/handlers"

	"github.com/gorilla/mux"
)

func RegisterRequestRoutes(router *mux.Router) {
	r := router.PathPrefix("/api/request/").Subrouter()

	r.HandleFunc("/{id}", handlers.GetRequestById).Methods("GET")

	r.HandleFunc("/method/{method_name}", handlers.GetRequestByMethod).Methods("GET")
}
