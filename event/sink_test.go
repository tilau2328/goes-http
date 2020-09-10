package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tilau2328/goes"
	"github.com/tilau2328/goes/core/event"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ExpectedEventResult = "event"

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/test", testMock)
	srv := httptest.NewServer(handler)
	return srv
}

func testMock(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	request := TestEvent{}
	err = json.Unmarshal(body, &request)
	if err != nil || request.Value != ExpectedEventResult {
		return
	}
	_, _ = res.Write([]byte(ExpectedHandlerResult))
}

func TestNewSink(t *testing.T) {
	sink := NewSink("", &http.Client{}, nil)
	if sink == nil {
		t.Errorf("failed to create event sink")
	}
}

func TestSink_Handle(t *testing.T) {
	srv := serverMock()
	defer srv.Close()
	sink := NewSink(srv.URL+"/test", &http.Client{}, "")
	c := event.NewEvent(uuid.New(), uuid.New(), TestEvent{ExpectedEventResult})
	response, err := sink.Handle(c)
	if err != nil {
		t.Error(err)
	}
	if response != ExpectedHandlerResult {
		t.Errorf("expected result to be %s but was %s", ExpectedHandlerResult, response)
	}
}

func TestSink_Register(t *testing.T) {
	bus := event.NewBus()
	sink := NewSink("", &http.Client{}, nil)
	sink.Register(bus, (*TestEvent)(nil))
	handler := bus.Handler(goes.MessageType((*TestEvent)(nil)))
	if handler != sink {
		t.Errorf("expected handler to be %T but was %T", sink, handler)
	}
}
