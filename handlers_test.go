package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Testing just handlers without router
// func TestCheckHealth(t *testing.T) {
// 	tt := []struct {
// 		name   string
// 		url    string
// 		exp    []byte
// 		status int
// 	}{
// 		{name: "first test", url: "/health", exp: []byte(`{"Message":"OK"}`), status: http.StatusOK},
// 	}
//
// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Create the new request
// 			req, err := http.NewRequest("GET", "localhost"+app.Conf.Addr+tc.url, nil)
// 			if err != nil {
// 				t.Fatalf("could not create request")
// 			}
// 			// Create new recorder, receive and save the handler response
// 			rec := httptest.NewRecorder()
// 			app.CheckHealth(rec, req)
// 			res := rec.Result()
// 			defer res.Body.Close()
// 			// Read the response
// 			got, err := ioutil.ReadAll(res.Body)
// 			if err != nil {
// 				t.Fatalf("cannot read the response body")
// 			}
// 			// Check Status
// 			if res.StatusCode != tc.status {
// 				t.Errorf("expected status %v; got %v", tc.status, res.StatusCode)
// 			}
// 			// Compare expected and actual body content
// 			if !bytes.Equal(got, tc.exp) {
// 				t.Errorf("expected response to be %v, got %v", string(tc.exp), string(got))
// 			}
// 		})
// 	}
// }

func TestCheckHealth(t *testing.T) {
	tt := []struct {
		name   string
		method string
		url    string
		body   []byte
		exp    interface{}
		status int
	}{
		{name: "health get", method: "GET", url: "/health", exp: map[string]interface{}{"Message": "OK"}, status: http.StatusOK},
		{name: "health post", method: "POST", url: "/health", status: http.StatusMethodNotAllowed},
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
				t.Fatalf("could not create request: ", err)
			}
			// Make the request
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("could not send %v request: %v", tc.method, err)
			}
			// Check Status
			if res.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, res.StatusCode)
			}
			// Read the response Body to got value
			var got interface{}
			err = json.NewDecoder(res.Body).Decode(&got)
			defer res.Body.Close()
			switch {
			case err == io.EOF:
			case err != nil:
				t.Fatalf("could not unmarshal data: ", err)
			}
			// Compare expected with received
			if !reflect.DeepEqual(tc.exp, got) {
				t.Errorf("expected response to be %v, got %v", tc.exp, got)
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	tt := []struct {
		name   string
		method string
		url    string
		body   []byte
		exp    map[string]interface{}
		status int
	}{
		{name: "first post",
			method: "POST",
			url:    "/posts",
			body:   []byte(`{"Name": "Post1", "Body": "Body1", "Tags": [{"Name": "tag1"},{"Name": "tag2"}]}`),
			// exp:    Post{ID: 1, Name: "Post1", Body: "Body1", Tags: []Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}},
			exp:    map[string]interface{}{"ID": float64(1), "Name": "Post1", "Body": "Body1", "Tags": []interface{}{map[string]interface{}{"ID": float64(1), "Name": "tag1"}, map[string]interface{}{"ID": float64(2), "Name": "tag2"}}},
			status: http.StatusOK,
		},
		{name: "with the same post name",
			method: "POST",
			url:    "/posts",
			body:   []byte(`{"Name": "Post1", "Body": "Body1", "Tags": [{"Name": "tag3"},{"Name": "tag4"}]}`),
			exp:    map[string]interface{}{"Message": "Failed. This post name already exists"},
			status: http.StatusBadRequest,
		},
		{name: "with the same tags",
			method: "POST",
			url:    "/posts",
			body:   []byte(`{"Name": "Post2", "Body": "Body2", "Tags": [{"Name": "tag1"},{"Name": "tag2"}]}`),
			// exp:    Post{ID: 2, Name: "Post2", Body: "Body2", Tags: []Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}},
			exp:    map[string]interface{}{"ID": float64(2), "Name": "Post2", "Body": "Body2", "Tags": []interface{}{map[string]interface{}{"ID": float64(1), "Name": "tag1"}, map[string]interface{}{"ID": float64(2), "Name": "tag2"}}},
			status: http.StatusOK,
		},
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
				t.Fatalf("could not create request: ", err)
			}
			// Make the request
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("could not send %v request: %v", tc.method, err)
			}
			// Read the response Body to post value
			// var post Post
			// err = json.NewDecoder(res.Body).Decode(&post)
			// defer res.Body.Close()
			// if err != nil {
			// 	t.Fatalf("could not unmarshal data: ", err)
			// }

			// Check Status
			if res.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, res.StatusCode)
			}

			// https://attilaolah.eu/2013/11/29/json-decoding-in-go/
			var rcv map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&rcv)
			defer res.Body.Close()
			if err != nil {
				t.Fatalf("could not unmarshal data: ", err)
			}
			// for i, y := range tc.exp {
			// 	fmt.Printf("%v -- %T: %v -- %T\n", i, i, y, y)
			// }
			// for i, y := range rcv {
			// 	fmt.Printf("%v -- %T: %v -- %T\n", i, i, y, y)
			// }

			// Compare expected with received
			if !reflect.DeepEqual(tc.exp, rcv) {
				t.Errorf("expected response to be %v, got %v", tc.exp, rcv)
			}
		})
	}
}
