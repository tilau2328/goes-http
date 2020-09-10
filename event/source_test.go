package event

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	utils "github.com/tilau2328/goes-http"
	"github.com/tilau2328/goes/core/event"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ExpectedHandlerResult = "test"

type TestEvent struct {
	Value string `json:"value"`
}
type TestHandler struct{}

func (*TestHandler) Handle(event.IEvent) (interface{}, error) {
	return ExpectedHandlerResult, nil
}

func TestNewSource(t *testing.T) {
	bus := event.NewBus()
	source := NewSource(bus, nil, func(interface{}, *http.Request) event.IEvent { return nil })
	if source == nil {
		t.Errorf("failed to create event source")
	}
}

func TestSource_Handle(t *testing.T) {
	bus := event.NewBus()
	message := &TestEvent{}
	aggregateId := uuid.New()
	source := NewSource(bus, (*TestEvent)(nil), func(body interface{}, r *http.Request) event.IEvent {
		return event.NewEvent(uuid.New(), utils.FirstId(r.RequestURI), body)
	})
	bus.RegisterHandler((*TestEvent)(nil), &TestHandler{})
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
