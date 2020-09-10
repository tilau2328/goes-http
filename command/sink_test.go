package command

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tilau2328/goes/core/command"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ExpectedCommandResult = "command"

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
	request := &TestCommand{}
	err = json.Unmarshal(body, request)
	if err != nil || request.Value != ExpectedCommandResult {
		return
	}
	_, _ = res.Write([]byte(ExpectedHandlerResult))
}

func TestNewSink(t *testing.T) {
	sink := NewSink(&http.Client{}, "", nil, func(body interface{}, response *http.Response) (interface{}, error) { return nil, nil })
	if sink == nil {
		t.Errorf("failed to create command sink")
	}
}

func TestSink_Handle(t *testing.T) {
	srv := serverMock()
	defer srv.Close()
	var res string
	sink := NewSink(&http.Client{}, srv.URL+"/test", res, func(body interface{}, response *http.Response) (interface{}, error) {
		return body, nil
	})
	aggregateId := uuid.New()
	message := &TestCommand{ExpectedCommandResult}
	c := command.NewCommand(uuid.New(), aggregateId, message)
	response, err := sink.Handle(c)
	if err != nil {
		t.Error(err)
	}
	if response != ExpectedHandlerResult {
		t.Errorf("expected result to be %s but was %s", ExpectedHandlerResult, response)
	}
}
