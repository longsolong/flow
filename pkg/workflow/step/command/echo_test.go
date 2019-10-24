package command

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEchoCommandType(t *testing.T) {
	echo := NewEchoCommand(uuid.New())
	assert.Equal(t, "github.com/longsolong/flow/pkg/workflow/step/command.EchoCommand", echo.Id().Type)


	/*
		expected := "application/json"
		t.Errorf("request header contains wrong content type: got %v want %v",
			"application/json", expected)
	*/
}


func TestEchoCommandRun(t *testing.T) {

}


