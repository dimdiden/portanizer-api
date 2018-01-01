package main

import (
	"encoding/json"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	body := map[string]string{"Status": "OK"}
	output, err := json.Marshal(&body)
	if err != nil {
		ErrorWithJSON(w, "Can not marshal output", http.StatusInternalServerError)
		return
	}
	ResponseWithJSON(w, output, http.StatusOK)
}
