package command

import (
	"bytes"
	"encoding/json"
	"github.com/tilau2328/goes/core/command"
	"io/ioutil"
	"net/http"
)

type Sink struct {
	client   *http.Client
	url      string
	encoding string
	response *interface{}
	cb       func(interface{}, *http.Response) error
}

func NewSink(
	client *http.Client,
	url string,
	encoding string,
	response *interface{},
	cb func(interface{}, *http.Response) error,
) *Sink {
	return &Sink{client, url, encoding, response, cb}
}

func (s *Sink) Handle(command command.ICommand) error {
	var body []byte
	var err error
	body, err = json.Marshal(command)
	if err != nil {
		return err
	}
	var res *http.Response
	res, err = http.Post(s.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if s.response == nil {
		return nil
	}
	err = json.Unmarshal(body, s.response)
	return s.cb(&s.response, res)
}
