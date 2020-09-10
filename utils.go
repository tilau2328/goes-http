package goes_http

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/tilau2328/goes"
	"net/http"
)

func Ids(uri string) []string {
	return goes.Regex("(\\b[0-9a-f]{8}\\b-([0-9a-f]{4}-){3}\\b[0-9a-f]{12}\\b)", uri)
}

func FirstId(uri string) uuid.UUID {
	return uuid.MustParse(Ids(uri)[0])
}

func RespondAndLogError(err error, w http.ResponseWriter) {
	if err != nil {
		fmt.Printf("error when writting response: %s\n", err.Error())
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("error when writting response: %s\n", err.Error())
		}
		return
	}
}
