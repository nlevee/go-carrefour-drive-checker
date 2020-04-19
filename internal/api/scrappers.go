package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nlevee/go-carrefour-drive-checker/pkg/carrefour"
)

// AddScrapperRoutes populate router
func addScrapperRoutes(r *mux.Router) {
	// ajoute un scrapper sur le store : storeid
	r.HandleFunc("/scrappers/{storeid}", addScrapper).Methods(http.MethodPut)
	// récupère l'état du scrapper sur le store : storeid
	r.HandleFunc("/scrappers/{storeid}", getScrapperState).Methods(http.MethodGet)
}

// GetScrapperState récuperation d'un état de scrapper
func getScrapperState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	log.Printf("storeId: %v\n", vars["storeid"])

	json.NewEncoder(w).Encode(carrefour.GetDriveState(vars["storeid"]))
}

// AddScrapper ajoute un scrapper
func addScrapper(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	vars := mux.Vars(r)
	log.Printf("storeId: %v\n", vars["storeid"])

	go carrefour.NewDriveHandler(vars["storeid"])
}
