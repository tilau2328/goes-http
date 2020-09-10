package query

import (
	"bytes"
	"encoding/json"
	"github.com/tilau2328/goes"
	"github.com/tilau2328/goes/core/query"
	"io/ioutil"
	"net/http"
)

type Sink struct {
	url      string
	client   *http.Client
	response interface{}
}

func NewSink(
	url string,
	client *http.Client,
	response interface{},
) *Sink {
	return &Sink{url, client, response}
}

func (s *Sink) Handle(c query.IQuery) (interface{}, error) {
	var res *http.Response
	response := s.response
	body, err := json.Marshal(c.Message())
	if err != nil {
		return nil, err
	}
	res, err = s.client.Post(s.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil || response == nil {
		return response, err
	}
	if goes.MessageType(response) == "string" {
		response = string(body)
	} else {
		err = json.Unmarshal(body, response)
	}
	return response, err
}

func (s *Sink) Register(bus query.IQueryBus, c interface{}) {
	bus.RegisterHandler(c, s)
}
