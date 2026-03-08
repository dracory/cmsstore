package shared

import "net/http"

func MediaManagerURL(r *http.Request) string {
	value := r.Context().Value(KeyMediaManagerURL)

	if value == nil {
		return ""
	}

	return value.(string)
}
