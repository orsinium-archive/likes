package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// SiteListHandler returns list of available sites
func SiteListHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/sites", SiteListHandler).Methods("GET")
	return r
}

func main() {
	r := getRouter()
	log.Fatal(http.ListenAndServe(":8000", r))
}
