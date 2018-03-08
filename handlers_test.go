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
	// Initiate server
	srv := httptest.NewServer(app.R)
	defer srv.Close()
	// Interate over the tasting table
	for _, tc := range TThealth {
		t.Run(tc.name, func(t *testing.T) {
			// Make the request and get the response
			res, err := makeRequest(srv, tc)
			if err != nil {
				t.Fatalf("could not create request: ", err)
			}
			// Check Status
			if res.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, res.StatusCode)
			}
			// Read the response Body to got value
			var rcv map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&rcv)
			defer res.Body.Close()
			switch {
			case err == io.EOF:
			case err != nil:
				t.Fatalf("could not unmarshal data: ", err)
			}
			// Compare expected with received
			if !reflect.DeepEqual(tc.exp, rcv) {
				t.Errorf("expected response to be %v, got %v", tc.exp, rcv)
			}
		})
	}
}

func TestPost(t *testing.T) {
	// Initiate server
	srv := httptest.NewServer(app.R)
	defer srv.Close()
	// Interate over the tasting table
	for _, tc := range TTpost {
		t.Run(tc.name, func(t *testing.T) {
			// Make the request and get the response
			res, err := makeRequest(srv, tc)
			if err != nil {
				t.Fatalf("could not create request: ", err)
			}
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
			// Compare expected with received
			if !reflect.DeepEqual(tc.exp, rcv) {
				t.Errorf("expected response to be %v, got %v", tc.exp, rcv)
			}
		})
	}
}

func makeRequest(srv *httptest.Server, tc Tc) (*http.Response, error) {
	// Prepare the request
	body := bytes.NewBuffer(tc.body)
	req, err := http.NewRequest(tc.method, srv.URL+tc.url, body)
	if err != nil {
		return nil, err
	}
	// Make the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Test case struct
type Tc struct {
	name   string
	method string
	url    string
	body   []byte
	exp    map[string]interface{}
	status int
}

// Set of unit tests for requests to health url
var TThealth = []Tc{
	{name: "health get",
		method: "GET",
		url:    "/health",
		exp:    map[string]interface{}{"Message": "OK"},
		status: http.StatusOK,
	},
	{name: "health post",
		method: "POST",
		url:    "/health",
		status: http.StatusMethodNotAllowed,
	},
}

// Set of unit tests for requests to post url
var TTpost = []Tc{
	// Create section
	{name: "first post",
		method: "POST",
		url:    "/posts",
		body:   []byte(`{"Name": "Post1", "Body": "Body1", "Tags": [{"Name": "tag1"},{"Name": "tag2"}]}`),
		exp:    map[string]interface{}{"ID": float64(1), "Name": "Post1", "Body": "Body1", "Tags": []interface{}{map[string]interface{}{"ID": float64(1), "Name": "tag1"}, map[string]interface{}{"ID": float64(2), "Name": "tag2"}}},
		status: http.StatusOK,
	},
	{name: "create with existing name",
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
		exp:    map[string]interface{}{"ID": float64(2), "Name": "Post2", "Body": "Body2", "Tags": []interface{}{map[string]interface{}{"ID": float64(1), "Name": "tag1"}, map[string]interface{}{"ID": float64(2), "Name": "tag2"}}},
		status: http.StatusOK,
	},
	{name: "create with with empty tags",
		method: "POST",
		url:    "/posts",
		body:   []byte(`{"Name": "Post3", "Body": "Body3"}`),
		exp:    map[string]interface{}{"ID": float64(3), "Name": "Post3", "Body": "Body3", "Tags": []interface{}{}},
		status: http.StatusOK,
	},
	{name: "create with empty post name",
		method: "POST",
		url:    "/posts",
		body:   []byte(`{"Name": "", "Body": ""}`),
		exp:    map[string]interface{}{"Message": "Failed. Name field is empty"},
		status: http.StatusBadRequest,
	},
	// Update section
	{name: "update post name",
		method: "PATCH",
		url:    "/posts/1",
		body:   []byte(`{"Name": "Post4", "Body": "Body4", "Tags": [{"Name": "tag1"},{"Name": "tag2"}]}`),
		exp:    map[string]interface{}{"ID": float64(1), "Name": "Post4", "Body": "Body4", "Tags": []interface{}{map[string]interface{}{"ID": float64(1), "Name": "tag1"}, map[string]interface{}{"ID": float64(2), "Name": "tag2"}}},
		status: http.StatusOK,
	},
	{name: "update post's tags",
		method: "PATCH",
		url:    "/posts/1",
		body:   []byte(`{"Name": "Post4", "Body": "Body4", "Tags": [{"Name": "tag1"},{"Name": "tag3"}]}`),
		exp:    map[string]interface{}{"ID": float64(1), "Name": "Post4", "Body": "Body4", "Tags": []interface{}{map[string]interface{}{"ID": float64(1), "Name": "tag1"}, map[string]interface{}{"ID": float64(3), "Name": "tag3"}}},
		status: http.StatusOK,
	},
	{name: "update post with empty tags",
		method: "PATCH",
		url:    "/posts/1",
		body:   []byte(`{"Name": "Post4", "Body": "Body4"}`),
		exp:    map[string]interface{}{"ID": float64(1), "Name": "Post4", "Body": "Body4", "Tags": []interface{}{}},
		status: http.StatusOK,
	},
	{name: "update with existing name",
		method: "PATCH",
		url:    "/posts/1",
		body:   []byte(`{"Name": "Post2", "Body": "Body2", "Tags": [{"Name": "tag1"},{"Name": "tag3"}]}`),
		exp:    map[string]interface{}{"Message": "Failed. This post name already exists"},
		status: http.StatusBadRequest,
	},
}
