package page_update

import (
	"net/http"

	"github.com/dracory/req"
)

func reqGetString(r *http.Request, key string) string {
	return req.GetStringTrimmed(r, key)
}
