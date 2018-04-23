package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ResponseWithJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	switch data := data.(type) {
	case string:
		str := fmt.Sprintf("{\"Message\": \"%s\"}", data)
		res := []byte(str)
		w.Write(res)
	default:
		res, err := json.Marshal(data)
		if err != nil {
			ResponseWithJSON(w, "Can not marshal output", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	}
}
