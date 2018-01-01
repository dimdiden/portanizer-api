package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	conf := GetConf("conf.json")

	r := mux.NewRouter()

	r.HandleFunc("/health", Health).Methods("GET")

	http.ListenAndServe(conf.Addr, r)
}
