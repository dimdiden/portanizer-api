package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCheckHealth(t *testing.T) {
	tt := []struct {
		name   string
		url    string
		exp    []byte
		status int
	}{
		{name: "first test", url: "/health", exp: []byte(`{"Message":"OK"}`), status: http.StatusOK},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Create the new request
			req, err := http.NewRequest("GET", "localhost"+app.Conf.Addr+tc.url, nil)
			if err != nil {
				t.Fatalf("could not create request")
			}
			// Create new recorder, receive and save the handler response
			rec := httptest.NewRecorder()
			app.CheckHealth(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			// Read the response
			got, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("cannot read the response body")
			}
			// Check Status
			if res.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, res.StatusCode)
			}
			// Compare expected and actual body content
			if !bytes.Equal(got, tc.exp) {
				t.Errorf("expected response to be %v, got %v", string(tc.exp), string(got))
			}
		})
	}
}

func TestRouting(t *testing.T) {
	tt := []struct {
		name   string
		method string
		url    string
		body   []byte
		// exp    []byte
		exp    map[string]interface{}
		status int
	}{
		// {name: "health route", method: "GET", url: "/health", body: []byte(""), exp: []byte(`{"Message":"OK"}`), status: http.StatusOK},
		{name: "health route", method: "GET", url: "/health", body: []byte(""), exp: map[string]interface{}{"Message": "OK"}, status: http.StatusOK},
		// {name: "create post", method: "POST", url: "/posts", body: []byte(`{"Name": "Hui","Body": "Testtest","Tags": [{"Name": "mysql"},{"Name": "ntp"}]}`), exp: []byte(""), status: http.StatusOK},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Initiate server
			srv := httptest.NewServer(app.R)
			defer srv.Close()
			// Prepare the request
			body := bytes.NewBuffer(tc.body)
			req, err := http.NewRequest(tc.method, srv.URL+tc.url, body)
			if err != nil {
				t.Fatalf("could not create request")
			}
			client := &http.Client{}
			// Make the request
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}
			// Check status code
			if res.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, res.StatusCode)
			}
			// // Read the response
			// got, err := ioutil.ReadAll(res.Body)
			// if err != nil {
			// 	t.Fatalf("cannot read the response body")
			// }
			// defer res.Body.Close()
			// // Compare expected and actual body content
			// if !bytes.Equal(got, tc.exp) {
			// 	t.Errorf("expected response to be %v, got %v", string(tc.exp), string(got))
			// }
			decoder := json.NewDecoder(res.Body)
			defer res.Body.Close()
			var got interface{}
			err = decoder.Decode(&got)
			// err := json.Unmarshal(b, &got)
			if err != nil {
				t.Fatalf("could not unmarshal the data")
			}
			// ????
			m := got.(map[string]interface{})
			// fmt.Println(m["Message"])
			if !reflect.DeepEqual(tc.exp, m) {
				t.Errorf("expected response to be %v, got %v", tc.exp, m)
			}
			// fmt.Println(eq)
		})
	}
}
