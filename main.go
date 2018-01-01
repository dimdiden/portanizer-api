package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/health", Health).Methods("GET")

	http.ListenAndServe(":8080", r)
}

func Health(w http.ResponseWriter, r *http.Request) {
	body := map[string]string{"Status": "OK"}
	output, err := json.Marshal(&body)
	if err != nil {
		ErrorWithJSON(w, "Can not marshal output", http.StatusInternalServerError)
		return
	}
	ResponseWithJSON(w, output, http.StatusOK)
}
