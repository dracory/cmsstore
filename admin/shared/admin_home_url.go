package shared

import "net/http"

func AdminHomeURL(r *http.Request) string {
	value := r.Context().Value(KeyAdminHomeURL)

	if value == nil {
		return ""
	}

	return value.(string)
}
