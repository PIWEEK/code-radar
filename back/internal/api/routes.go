package api

import (
	"net/http"
	"github.com/gorilla/mux"
)

var FileSystem http.FileSystem = nil

func CreateRoutes() *mux.Router {
	routes := mux.NewRouter().StrictSlash(false)

	routes.HandleFunc("/files", RetrieveFiles)

	if FileSystem != nil {
		routes.PathPrefix("/").Handler(http.FileServer(FileSystem))
	}

	return routes;
}
