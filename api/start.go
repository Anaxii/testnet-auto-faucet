package api

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StartAPI() {

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/ws", getWS).Methods("GET")

	r.Use(mux.CORSMethodMiddleware(r))
	log.Fatal(http.ListenAndServe(":8080", r))
}
