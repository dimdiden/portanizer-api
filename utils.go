package main

import (
	"encoding/json"
	"net/http"
)

func ResponseWithJSON(w http.ResponseWriter, data interface{}, code int) {
	out, err := json.Marshal(data)
	if err != nil {
		ErrorWithJSON(w, "Can not marshal output", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(out)
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	m := map[string]string{"Message": message}
	out, _ := json.Marshal(&m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(out)
}
