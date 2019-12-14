// +build integration

package flow

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunFlowHandler(t *testing.T) {
	method := "POST"
	uri := "/api/single_processor/flows/run"
	payload := `{
		"primaryRequestArgs": {
			"name": "ping",
			"version": 1
		},
		"requestArgs": {
			"hostname": "www.baidu.com",
			"timeout":  10,
			"interval": 1,
			"count":    1
		},
		"requestTags": [
			{"name": "aa", "value": "bb"}
		]
	}`

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(method, uri, strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RunFlowHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	expected = "application/json; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("handler returned wrong content type header: got %v want %v",
			contentType, expected)
	}
}
