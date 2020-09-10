package command

import (
	"encoding/json"
	"github.com/tilau2328/goes/core/command"
	"io/ioutil"
	"net/http"
)

type Source struct {
	bus      command.ICommandBus
	request  interface{}
	template func(interface{}, *http.Request) command.ICommand
}

func NewSource(
	bus command.ICommandBus,
	request interface{},
	template func(interface{}, *http.Request) command.ICommand,
) *Source {
	return &Source{bus, request, template}
}

func (s *Source) Handle(w http.ResponseWriter, req *http.Request, request interface{}) (interface{}, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, request)
	if err != nil {
		return nil, err
	}
	var result interface{}
	result, err = s.bus.Handle(s.template(request, req))
	body, err = json.Marshal(result)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(body)
	return result, err
}
