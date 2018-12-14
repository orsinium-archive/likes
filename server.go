package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/recoilme/slowpoke"
)

const sitesFile = ".storage/sites.slowpoke"

// SiteListHandler returns list of available sites
func SiteListHandler(w http.ResponseWriter, r *http.Request) {
	keys, _ := slowpoke.Keys(sitesFile, nil, 0, 0, true)
	sep := []byte("\n")
	w.Write(bytes.Join(keys, sep))
}

// SiteAddHandler adds new site
func SiteAddHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slowpoke.Set(sitesFile, []byte(vars["site"]), []byte(""))
	w.Write([]byte("OK!"))
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/sites", SiteListHandler).Methods("GET")
	r.HandleFunc("/sites/{site:[a-zA-Z\\.]+}", SiteAddHandler).Methods("POST")
	return r
}

func main() {
	defer slowpoke.CloseAll()
	r := getRouter()
	log.Fatal(http.ListenAndServe(":8000", r))
}
