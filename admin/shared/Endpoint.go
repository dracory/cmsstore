package shared

import "net/http"

func Endpoint(r *http.Request) string {
	value := r.Context().Value(KeyEndpoint)

	if value == nil {
		return ""
	}

	return value.(string)
}
