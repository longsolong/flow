package flow

import (
	"github.com/faceair/jio"
	"github.com/longsolong/flow/dev/workflows"
	"io/ioutil"
	"net/http"
)

// RunFlowValidator ...
var RunFlowValidator = jio.Object().Keys(jio.K{
	"primaryRequestArgs": jio.Object().Keys(jio.K{
		"namespace":    jio.String().Required(),
		"name":    jio.String().Required(),
		"version": jio.Number().Integer().Required(),
	}),
	"requestArgs": jio.Object().Required(),
	"requestTags": jio.Array().Required(),
})

// RunFlowHandler ...
func RunFlowHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jio.DefaultErrorHandler(w, r, err)
		return
	}
	data, err := jio.ValidateJSON(&body, RunFlowValidator)
	if err != nil {
		jio.DefaultErrorHandler(w, r, err)
		return
	}
	namespace := data["primaryRequestArgs"].(map[string]interface{})["namespace"]
	name := data["primaryRequestArgs"].(map[string]interface{})["name"]
	version := data["primaryRequestArgs"].(map[string]interface{})["version"]
	_, err = workflows.SingleProcessorFactory.Make(namespace.(string), name.(string), int(version.(float64)), body)
	if err != nil {
		jio.DefaultErrorHandler(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return
}
