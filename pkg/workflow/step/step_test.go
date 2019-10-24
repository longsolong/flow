package step

import (
	"fmt"
	"reflect"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"
)

type EchoCommand struct {
	*ShellCommand
}

func NewEchoCommand(uuid uuid.UUID) *EchoCommand {
	empty := reflect.TypeOf((*EchoCommand)(nil)).Elem()
	return &EchoCommand{
		ShellCommand: &ShellCommand{
			Atom: &Atom{
				id: Id{Uuid: uuid, Type: fmt.Sprintf("%s.%s", empty.PkgPath(), empty.Name())},
			},
		},
	}
}


func TestStepType(t *testing.T) {
	var echo Step = NewEchoCommand(uuid.New())
	assert.Equal(t, echo.Id().Type, "github.com/longsolong/flow/pkg/workflow/step.EchoCommand")


	/*
	expected := "application/json"
	t.Errorf("request header contains wrong content type: got %v want %v",
		"application/json", expected)
	*/
}


func TestEchoCommandRun(t *testing.T) {

}

