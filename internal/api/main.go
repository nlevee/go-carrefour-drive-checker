package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// StartServer demarrage du server
func StartServer(host string, port string) {
	r := mux.NewRouter()
	addStoreRoutes(r)
	addScrapperRoutes(r)
	log.Fatal(http.ListenAndServe(host+":"+port, r))
}