package shared

import "net/http"

func Endpoint(r *http.Request) string {
	return r.Context().Value(KeyEndpoint).(string)
}
