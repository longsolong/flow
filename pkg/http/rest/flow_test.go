// +build integration

package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/longsolong/flow/pkg/infra"
)

func TestRunFlowHandler(t *testing.T) {
	method := "POST"
	uri := "/api/single_processor/flows/run"
	payload := `{
		"primaryRequestArgs": {
			"namespace": "examples",
			"name": "number_guess",
			"version": 1
		},
		"requestArgs": {
			"secret": 1
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

	logger, err := infra.CreateLogger(0)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SingleProcessorFlowHandler{logger: logger}.Run())

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"Atom":{"ID":"","ExpansionDigest":"","NumberGuessParameter":{"Secret":1,"Low":1,"High":1}},"State":0,"StateText":"SUCCESS"}]`
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
