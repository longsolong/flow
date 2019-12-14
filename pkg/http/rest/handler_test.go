package rest

import (
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/longsolong/flow/test"
)

func TestCreateHandler(t *testing.T) {
	logger, _ := test.LoggerMock()
	handler := CreateHandler(logger)

	expectedRouterType := reflect.TypeOf(chi.NewRouter())
	if r := reflect.TypeOf(handler.GetRouter()); r != expectedRouterType {
		t.Errorf("handler has wrong type of router: got %v want %v",
			r, expectedRouterType)
	}
}
