package shared

import "net/http"

func AdminHomeURL(r *http.Request) string {
	return r.Context().Value(KeyAdminHomeURL).(string)
}
