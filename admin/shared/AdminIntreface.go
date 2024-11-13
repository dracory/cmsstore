package shared

import "net/http"

type Admin interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
