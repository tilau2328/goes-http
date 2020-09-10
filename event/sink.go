package event

import (
	"bytes"
	"encoding/json"
	"github.com/tilau2328/goes"
	"github.com/tilau2328/goes/core/event"
	"io/ioutil"
	"net/http"
)

type Sink struct {
	client   *http.Client
	url      string
	response interface{}
	cb       func(interface{}, *http.Response) (interface{}, error)
}

func NewSink(
	client *http.Client,
	url string,
	response interface{},
	cb func(interface{}, *http.Response) (interface{}, error),
) *Sink {
	return &Sink{client, url, response, cb}
}

func (s *Sink) Handle(c event.IEvent) (interface{}, error) {
	var body []byte
	var err error
	body, err = json.Marshal(c.Message())
	if err != nil {
		return nil, err
	}
	var res *http.Response
	res, err = s.client.Post(s.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if s.response == nil {
		return nil, nil
	}
	if goes.MessageType(s.response) == "string" {
		s.response = string(body)
	} else {
		err = json.Unmarshal(body, s.response)
	}
	return s.cb(s.response, res)
}
