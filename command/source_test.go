package command

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	utils "github.com/tilau2328/goes-http"
	"github.com/tilau2328/goes/core/command"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ExpectedHandlerResult = "test"

type TestCommand struct {
	Value string `json:"Value"`
}
type TestHandler struct{}

func (*TestHandler) Handle(command.ICommand) (interface{}, error) {
	return ExpectedHandlerResult, nil
}

func TestNewSource(t *testing.T) {
	source := NewSource(command.NewBus(), nil, func(interface{}, *http.Request) command.ICommand { return nil })
	if source == nil {
		t.Errorf("failed to create command source")
	}
}

func TestSource_Handle(t *testing.T) {
	var result []byte
	bus := command.NewBus()
	message := &TestCommand{ExpectedCommandResult}
	aggregateId := uuid.New()
	source := NewSource(bus, message, func(body interface{}, r *http.Request) command.ICommand {
		return command.NewCommand(uuid.New(), utils.FirstId(r.RequestURI), body)
	})
	bus.RegisterHandler((*TestCommand)(nil), &TestHandler{})
	response := httptest.NewRecorder()
	b, err := json.Marshal(message)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("post", "/"+aggregateId.String()+"/test", bytes.NewReader(b))
	source.Handle(response, req)
	result, err = ioutil.ReadAll(response.Result().Body)
	if string(result) != "\""+ExpectedHandlerResult+"\"" {
		t.Errorf("expected handler response to be %s but was %s", "\""+ExpectedHandlerResult+"\"", string(result))
	}
}
