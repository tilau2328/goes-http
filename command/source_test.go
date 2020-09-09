package command

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tilau2328/goes/core/command"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ExpectedHandlerResult = "test"

type TestCommand struct{}
type TestHandler struct{}

func (*TestHandler) Handle(command.ICommand) (interface{}, error) {
	return ExpectedHandlerResult, nil
}

func TestNewSource(t *testing.T) {
	bus := command.NewBus()
	source := NewSource(bus, nil, func(interface{}, *http.Request) command.ICommand { return nil })
	if source == nil {
		t.Errorf("failed to create command source")
	}
}

func TestSource_Handle(t *testing.T) {
	bus := command.NewBus()
	message := &TestCommand{}
	aggregateId := uuid.New()
	source := NewSource(bus, (*TestCommand)(nil), func(body interface{}, r *http.Request) command.ICommand {
		return command.NewCommand(uuid.New(), FirstId(r.RequestURI), body)
	})
	bus.RegisterHandler((*TestCommand)(nil), &TestHandler{})
	response := httptest.NewRecorder()
	b, err := json.Marshal(message)
	if err != nil {
		t.Error(err)
	}
	var res interface{}
	req := httptest.NewRequest("post", "/"+aggregateId.String()+"/test", bytes.NewReader(b))
	res, err = source.Handle(response, req, message)
	if err != nil {
		t.Error(err)
	}
	if res != ExpectedHandlerResult {
		t.Errorf("expected response to be %s but was %s", ExpectedHandlerResult, res)
	}
	var result []byte
	result, err = ioutil.ReadAll(response.Result().Body)
	if string(result) != "\""+ExpectedHandlerResult+"\"" {
		t.Errorf("expected handler response to be %s but was %s", "\""+ExpectedHandlerResult+"\"", string(result))
	}
}
