package goes_http

import (
	"github.com/google/uuid"
	"github.com/tilau2328/goes"
)

func Ids(uri string) []string {
	return goes.Regex("(\\b[0-9a-f]{8}\\b-([0-9a-f]{4}-){3}\\b[0-9a-f]{12}\\b)", uri)
}

func FirstId(uri string) uuid.UUID {
	return uuid.MustParse(Ids(uri)[0])
}
