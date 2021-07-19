package api

import (
	"log"
	"net/http"
	"github.com/rs/cors"
)

func Start() error {
	routes := CreateRoutes()
	handler := cors.Default().Handler(routes)
	log.Println("Listening http://localhost:8000/files")
	return http.ListenAndServe(":8000", handler)
}
