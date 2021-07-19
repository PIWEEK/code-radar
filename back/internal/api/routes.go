package api

import (
	"github.com/gorilla/mux"
)

func CreateRoutes() *mux.Router {
	routes := mux.NewRouter().StrictSlash(false)

	routes.HandleFunc("/files", RetrieveFiles)

	return routes;
}
