package query

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	utils "github.com/tilau2328/goes-http"
	"github.com/tilau2328/goes/core/query"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ExpectedHandlerResult = "test"

type TestQuery struct {
	Value string `json:"value"`
}
type TestHandler struct{}

func (*TestHandler) Handle(query.IQuery) (interface{}, error) {
	return ExpectedHandlerResult, nil
}

func TestNewSource(t *testing.T) {
	bus := query.NewBus()
	source := NewSource(bus, nil, func(interface{}, *http.Request) query.IQuery { return nil })
	if source == nil {
		t.Errorf("failed to create query source")
	}
}

func TestSource_Handle(t *testing.T) {
	bus := query.NewBus()
	message := &TestQuery{}
	aggregateId := uuid.New()
	source := NewSource(bus, (*TestQuery)(nil), func(body interface{}, r *http.Request) query.IQuery {
		return query.NewQuery(uuid.New(), utils.FirstId(r.RequestURI), body)
	})
	bus.RegisterHandler((*TestQuery)(nil), &TestHandler{})
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
