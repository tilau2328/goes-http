package command

import (
	"encoding/json"
	HTTP "github.com/tilau2328/goes-http"
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

func (s *Source) Handle(w http.ResponseWriter, req *http.Request) {
	var response interface{}
	body, err := ioutil.ReadAll(req.Body)
	HTTP.RespondAndLogError(err, w)
	err = json.Unmarshal(body, s.request)
	HTTP.RespondAndLogError(err, w)
	response, err = s.bus.Handle(s.template(s.request, req))
	HTTP.RespondAndLogError(err, w)
	body, err = json.Marshal(response)
	HTTP.RespondAndLogError(err, w)
	_, err = w.Write(body)
	HTTP.RespondAndLogError(err, w)
}

func (s *Source) Register(pattern string) {
	http.HandleFunc(pattern, s.Handle)
}
