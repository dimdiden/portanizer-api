package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestCheckHealth (t *testing.T) {
	req, err := http.NewRequest("GET", "localhost:8080/health", nil)
	if err != nil {
		t.Fatalf("could not create request")
	}
	rec := httptest.NewRecorder()

	app.CheckHealth(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
	}
}
