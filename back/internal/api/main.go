package api

import (
	"log"
	"strconv"
	"net/http"
	"github.com/rs/cors"
)

func Start(port int) error {
	routes := CreateRoutes()
	handler := cors.Default().Handler(routes)
	strPort := strconv.Itoa(port)
	log.Println("Listening on: http://localhost:" + strPort)
	return http.ListenAndServe(":" + strPort, handler)
}
